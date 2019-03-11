## 2019-03-10
**ClientVerify接口传参**

>parse_url($params, PHP_URL_QUERY)

**query string的值**

| name | type | 
| ----- | ----- |  
| token | string | 
| user_id | string | 
| data | string|

**其中 token 值**

| name | type | 
| ----- | ----- | 
| token[0] | jsonData |
| token[1] | osdkConfId|
| token[2] | channelId|
| token[3] | time|

## 2019-03-11 code view

> ### EntranceController

* actionClientVerify ()
    * 进行二次验证的入口方法
    * 它首先拼接了 $_POST + $_GET ，前者若存在key，则以前面的为准，然后调用了 verify 方法
* verify ($putData, $ticket = 0)
    * 有两个参数，第二个参数 ticket: 1 表示返回信息，0 表示不返回。有另一个方法 actionServerVerify 调用 verify 时 ticket 传的为 0，但胖子说该方法已废弃不用了。
    * 由上面的传参可知，首先接受三个传入的值，分别为 $userId, $token, $data
    * 其中 $userId 是从渠道那边返回的，此处为二次验证需要
    * token 参数
        * `token 值不能为空，为空则报 status:11` 
        * token 值base64解开后，用 & 符号进行了拼接，分别为 
        * `token 解开后必须包含以下 4 个值，或者说代码里判断的标准是数量必须是 4 个，否则报 status:12`
        * 0: jsonData（为验证时需要从 supersdk 前端传过来的值，进行了 rsa 加密，此处需要进行 rsa 私钥解密）
            * jsonData 解开后的值是跟前端商量好的传参，`如果解密后为空，则报 status:13`
            * jsonData 解开后的值要传递给相应的`业务逻辑控制器`去接受，此处存为 $arrParams
        * 1: osdk_conf_id（为打包方案配置 id，即 sdk_conf 表的 id，sdk_conf 和 sdk_detail 表在打包方案进行一次打包方案生成的适合，会从 scheme_conf 和 scheme_detail 表中拷贝一份快照数据)，这个值在客户端里保存，在打包的适合一并打入包中
            * 在获取 loginSdkConf 信息的时候，通过该 osdk_conf_id 从 sdk_conf 表中连接了 sdk 和 sdk_detail 表来获取有用信息 osdkConfInfo
            * `标记1` 这里需要注意的是，这里获取的数据会在 redis 里保存一份 hash 数据，key 为 MGOS-LoginSDKConf-{$osdkConfId}，最终缓存用的是 hMset，hMset 的特性是 删除了 key 中的 field，再次 hMset 的时候，并不会将该 field 删除掉，此时需要手动清理 redis 里的数据，或者让前端重新出一次打包方案，因为重新出，sdk_conf 会生成一次新的，所以也就重新读库入缓存了。
            * `如果这里获取 osdkConfInfo 信息为空，则报 status:14`
            * 从上步 osdkConfInfo 中可以获得游戏信息(_game_id)，获取游戏信息，游戏在配置的时候，会自动生成 game_server_secret，此参数为必须要有的参数（合服的时候也要注意此参数），`如果该参数缺失或者为空，则报 status:17`
            * 上步获取 GameInfo 是通过 game 表的主键 game_id 来获取的，获取 game_id,game_name,debug_notify_url,release_notify_url,game_server_secret 信息后，hmset 到 redis 里，同样存在上面`标记1`的问题。
        * 2: channelId，胖子说我们不关心这个值 ，属于透传给游戏或者其他部门使用
        * 3: time 时间戳
        * 以上 4 个参数处理完毕后，将 is_ticket,user_id,channel_id一并打包到 $arrParams 里传给`业务逻辑控制器`
    * data 参数 
        * demo: `{"device_id":"862224033038327_38398a0387cd64d9_0","framework_version":"哆可梦(2.5.10)|哆可梦(2.5.10)","login_version":"哆可梦(20190104-7.0)","network":"WIFI","network_operators":"other","os":"android","osversion":"7.1.1","pay_version":"哆可梦(20190104-7.0)","phone_model":"PRO 6","plugin_version":"哆可梦(1.0.0)|哆可梦(1.0.0)","screen_resolution":"1920*1080","sdk_version":"4.0.5"}`
        * 从 demo 能看出，data 参数的值解开后，附加上 osdk_conf_id,game_id,timeline,user_id(若不为空，则前缀加上 account_system_id 和 _，如传过来的 uid=45234799,则为 0060015_45234799),ip,type=1，打包这些数据后 ，一并写入 flume 里(log_osdk_client_basic)
    * 在处理完入参后，得到需要的两个值 osdkConfInfo 和  arrParams 传入 dispatch 方法，开始进行路由跳转前的预处理
* dispatch ($osdkConfInfo, $action, $params = [])
    * 代理函数，处理路由转发逻辑
    * 3 个参数，第一个参数是 osdkConfInfo 信息，verify 步传入，action 是之后`业务逻辑控制器` 里处理具体业务逻辑的方法，此处登录为固定的 login （上步 verify 传入），第 3 个参数 params 是 verify 传入的客户端参数和一些拼接参数
    * `此步骤会有 3 种情况报  status:15`
        * `业务逻辑控制器` 未继承 ControllerApi
        * `业务逻辑控制器` 中没有 actionLogin 方法
        * `业务逻辑控制器` 不存在 （较为多见：比如传了金华 ，忘传多湖了；比如 sdk 配置信息，金华填写了，忘填（或者填错）多湖了等情况）
    * 在此步骤中，会在 `业务逻辑控制器` 传入 4 个属性
        * MR...Ctl->_params = $arrParams
        * MR...Ctl->_conf = $osdkConfInfo
        * MR...Ctl->_osdk_conf_id = $osdkConfInfo['_osdk_conf_id']
        * MR...Ctl->_game_server_secret = $osdkConfInfo['_game_server_secret']

> ### MR...Controller extends ControllerApi

* actionLogin ()
    * 入口一个 log "登录请求入口，info 打印 $this->params"
    * 获取 3 个参数， osdkConfId,user_id, gameServerSecret(is_ticket=1的时候获取，ClientVerify请求时为1)
    * `客户端透传的参数为空，则报 status:21`
    * url 此处为请求渠道的二次验证的地址，从文档中获取，加上透传的参数，进行 curl 请求，获取返回值 $rst
    * `如果上步的 $rst 的值为空，则认为渠道二次验证的响应出错，报 status:24`
    * `如果 $rst 不为空，并且符合渠道的验证规则，则通过 FunCom::loginEcho 将 登录验证成功输出 给客户端，否则为渠道验证失败，报 status:23`

> ### FunCom::loginEcho

* 入参
    1. @param string $userId 用户标识
    2. @param string $loginSdkName 登录SDK名称
    3. @param string $gameServerSecert 密钥
    4. @param integer $osdkConfId
    5. @param array $data 登录成功后第三方联运平台返回的扩展参数
    6. @param string $apiName 联运平台名称
    7. @param array $param 登录请求所有参数 直接传data
    8. @param array $osdkConfInfo osdk的配置信息
* 对于入参，`如果 8 中的用户体系id（前端小鱼或者油纸伞跟商务要） _account_system_id 为空，则报 status:15`
* 该方法也有个 gameServerSecret 是否为空的判断，在 EntranceCtroller->verify 中也有判断，所以该处的判断应该是防止他处的调用，哪处呢？(搜过没搜到)
* 假设上个步骤是含有 gameServerSecret 不为空的，则此步骤继续打包数组 $arrData（该数组将会返回给前端），包含 user_id,loginSdkName,ip,country,time,sign
* 上步的 ip 需要个 ip 库的存在，此处是从 redis 中 zset 的数据结构中，key=SP-IP,主要是为了获取国家信息，如“CN:AP:2”，如无则为“00”
* 出参
    * status： 1
    * msg：登录成功
    * userInfo: 用户信息
        * user_id
        * login_sdk_name
    * data: 登录渠道后返回的数据
    * osdk_ticket：上面打包的数组 $arrData, 进行 base64(json) 操作