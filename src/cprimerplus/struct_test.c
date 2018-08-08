#include <stdio.h>
#include <stdlib.h>

typedef struct Node {
	int data;
	LinkNode *next;
} LinkNode;

int main(void)
{
	LinkNode *lnode1, *lnode2;	
	lnode1 = (LinkNode *)malloc(sizeof(LinkNode *));
	lnode2 = (LinkNode *)malloc(sizeof(LinkNode *));

	lnode1->data = 123;
	lnode1->next = lnode2;

	lnode2->data = 456;
	lnode2->next = NULL;

	
	return 0;
}
