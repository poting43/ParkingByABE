#include <iostream>
#include <ctime>

//********* choose just one of these pairs **********
//#define MR_PAIRING_CP      // AES-80 security   
//#define AES_SECURITY 80

//#define MR_PAIRING_MNT	// AES-80 security
//#define AES_SECURITY 80

#define MR_PAIRING_BN    // AES-128 or AES-192 security
#define AES_SECURITY 128
//#define AES_SECURITY 192

//#define MR_PAIRING_KSS    // AES-192 security
//#define AES_SECURITY 192

//#define MR_PAIRING_BLS    // AES-256 security
//#define AES_SECURITY 256
//*********************************************

#include "pairing_3.h"


#define HASH_LEN 32
Big mysha(char *string)
{
    Big h;
    unsigned char s[HASH_LEN];
    int i,j; 
    sha256 sh;

    shs256_init(&sh);

    for (i=0;;i++)
    {
        if (string[i]==0) break;
        shs256_process(&sh,string[i]);
    }
    shs256_hash(&sh,(char *)s);
    h=1;
    for(i=0;i<32;i++)
    {
      h*=(unsigned int)s[i];
      h*=256;
    }
    
    return h;
}

int main()
{  
  clock_t kgen_start, kgen_finish;
  double kgen_duration;
	PFC pfc(AES_SECURITY);  // initialise pairing-friendly curve
  cout<<"initialise pairing-friendly curve"<<endl; 
	miracl *mip=get_mip();  // get handle on mip (Miracl Instance Pointer)
	Big order=pfc.order();  // get pairing-friendly group order

	time_t seed;            // crude randomisation
	time(&seed);
    irand((long)seed);
  cout<<"gan"<<endl;
  G1 p1;
  G2 p2;
  GT YA,YB;
  Big ab,b;
  pfc.random(p1);
  pfc.random(p2);
  pfc.random(ab);
  b=16;
  cout<<"b = "<<b<<endl;
  cout<<"b*10 = "<<b*10<<endl;
  cout<<"p1 = "<<p1.g<<endl;
  cout<<"p2 = "<<p2.g<<endl;
  Big p11,p12,p211,p212,p221,p222;
  p1.g.get(p11,p12);
  ZZn2 x,y;
  p2.g.get(x,y);
  x.get(p211,p212);
  y.get(p221,p222);
  cout<<"p11 = "<<p11<<endl;
  cout<<"p12 = "<<p12<<endl;
  cout<<"p211 = "<<p211<<endl;
  cout<<"p212 = "<<p212<<endl;
  cout<<"p221 = "<<p221<<endl;
  cout<<"p222 = "<<p222<<endl;
  char c[200];
  mip->IOBASE=10;
  c<<b;
  cout<<"c = "<<c<<endl;
  mip->IOBASE=16;
  char cb1[500];
  cb1<<p11;
  cout<<"cb1 = "<<cb1<<endl;
  char cb2[500];
  cb2<<p12;
  cout<<"cb2 = "<<cb2<<endl;
  strcat(cb1,cb2);
  cout<<"cb1 = "<<cb1<<endl;
  char c123[]="123";
  cout<<mysha(c123)<<endl;
  
  
  return 0;
}