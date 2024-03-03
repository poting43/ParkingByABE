/*
   Constant Size Ciphertexts in Threshold ABE - Javier Herranz, Fabien Laguillaumie, and Carla R`afols
   See https://www.lirmm.fr/~laguillaum/shortABEpkc10.pdf
   Section 3.1

   Compile with modules as specified below

   For MR_PAIRING_BN curve
   cl /O2 /GX time_tbcpabe.cpp bn_pair.cpp zzn12a.cpp ecn2.cpp zzn4.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp miracl.lib

   In linux: replace "cl /O2 /GX" to "g++", "miracl.lib" to "miracl.a"

   g++ time_sign.cpp bn_pair.cpp zzn12a.cpp ecn2.cpp zzn4.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp mrshs256.c miracl.a

   Test program
*/

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


int main()
{
  clock_t kgen_start, kgen_finish;
  double kgen_duration;
	PFC pfc(AES_SECURITY);  // initialise pairing-friendly curve
	miracl *mip=get_mip();  // get handle on mip (Miracl Instance Pointer)
	Big order=pfc.order();  // get pairing-friendly group order

	time_t seed;            // crude randomisation
	time(&seed);
    irand((long)seed);
    
// pidgen
  double START,END; 
  START = clock();
  Big x;
  pfc.random(x);
  G1 Q1;
  G2 Q2;
  pfc.random(Q1);
  pfc.random(Q2);
	mip->IOBASE=10;
	mip->IOBASE=16;
  G1 P1;
  P1=pfc.mult(Q1,x);
  Big ID;
  pfc.random(ID);
  G1 Z;
  Big r3;
  pfc.random(r3);
  Z=pfc.mult(Q1,r3);
  Big P2;
  pfc.start_hash();
  pfc.add_to_hash(pfc.mult(Z,x));
  P2=lxor(ID,pfc.finish_hash_to_aes_key());
  Big M;
  END = clock();
  cout << endl << "pidgen : " << (END - START) / CLOCKS_PER_SEC << "s" << endl;

// sign
  START = clock();
  G2 s;
  pfc.start_hash();
  pfc.add_to_hash(M);
  pfc.add_to_hash(P1);
  pfc.add_to_hash(P2);
  Big T;
  pfc.random(T);
  pfc.add_to_hash(T);
  Big hashm = pfc.finish_hash_to_aes_key();
  s=pfc.mult(Q2,x+modmult(r3,hashm,order));
  
  END = clock();
  cout << endl << "Sign : " << (END - START) / CLOCKS_PER_SEC << "s" << endl;
// Trace
  START = clock();
  Big NID;
  pfc.start_hash();
  pfc.add_to_hash(pfc.mult(P1,r3));
  NID=lxor(P2,pfc.finish_hash_to_aes_key());
  if(ID==NID)
    cout<<"Trace correct"<<endl;
    
  END = clock();
  cout << endl << "Trace : " << (END - START) / CLOCKS_PER_SEC << "s" << endl;
// Verify
  START = clock();
  GT left;
  GT right;
  left=pfc.pairing(s,Q1);
  right=pfc.pairing(Q2,P1+pfc.mult(Z,hashm));
  if(left==right)
    cout<<"Verify True"<<endl;
  else cout<<"Verify False"<<endl;
  
  
  END = clock();
  cout << endl << "Verify : " << (END - START) / CLOCKS_PER_SEC << "s" << endl;


  return 0;
}
