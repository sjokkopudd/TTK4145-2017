#include <pthread.h>
#include <stdio.h>

int i = 0;

void *inc_i(void *i_void) {
	
	int *i = (int *)i_void;
	
	for(int j = 0; j < 1000000; j++){
		(*i)++;
	}

	return NULL;
}

void *dec_i(void *i_void) {
	
	int *i = (int *)i_void;
	
	for(int j = 0; j < 1000000; j++){
		(*i)--;
	}

	return NULL;
}

int main(){

	pthread_t inc_thread;
	pthread_create(&inc_thread, NULL, inc_i, &i);
	
	pthread_t dec_thread;
	pthread_create(&dec_thread, NULL, dec_i, &i);

	printf("i: %i \n", i);

	pthread_join(inc_thread, NULL);
	pthread_join(dec_thread, NULL);

	return 1; 

	

}
