/* test program: should produce digest  

248d6a61 d20638b8 e5c02693 0c3e6039 a33ce459 64ff2167 f6ecedd4 19db06c1
*/

#include <stdio.h>
#include "miracl.h"
#include <iostream>
#include <cstring>

using namespace std;

char test[]="abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq";

int main()
{
    char hash[32];
    char output[64]={0};
    int i;
    sha256 sh;
    shs256_init(&sh);
    for (i=0;test[i]!=0;i++) shs256_process(&sh,test[i]);
    shs256_hash(&sh,hash);    
    for (i=0;i<32;i++) 
    {
      char c[2];
      sprintf(c, "%02x",(unsigned char)hash[i]);
      strcat(output,c);
    }
    cout<<output<<endl;
    return 0;
}

