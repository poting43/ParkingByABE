/*
   Constant Size Ciphertexts in Threshold ABE - Javier Herranz, Fabien Laguillaumie, and Carla R`afols
   See https://www.lirmm.fr/~laguillaum/shortABEpkc10.pdf
   Section 3.1

   Compile with modules as specified below

   For MR_PAIRING_CP curve
   cl /O2 /GX myabe.cpp cp_pair.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp miracl.lib

   For MR_PAIRING_MNT curve
   cl /O2 /GX myabe.cpp mnt_pair.cpp zzn6a.cpp ecn3.cpp zzn3.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp miracl.lib

   For MR_PAIRING_BN curve
   cl /O2 /GX myabe.cpp bn_pair.cpp zzn12a.cpp ecn2.cpp zzn4.cpp zzn2.cpp big.cpp zzn.cpp ecn.cpp miracl.lib

   For MR_PAIRING_KSS curve
   cl /O2 /GX myabe.cpp kss_pair.cpp zzn18.cpp zzn6.cpp ecn3.cpp zzn3.cpp big.cpp zzn.cpp ecn.cpp miracl.lib

   For MR_PAIRING_BLS curve
   cl /O2 /GX myabe.cpp bls_pair.cpp zzn24.cpp zzn8.cpp zzn4.cpp zzn2.cpp ecn4.cpp big.cpp zzn.cpp ecn.cpp miracl.lib

   In linux: replace "cl /O2 /GX" to "g++", "miracl.lib" to "miracl.a"

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

//
//  Note that in this case Bobs attributes are "close enough" to Alices
//  so that he can decrypt
//

#define NATTR  20 // Universe of attributes
#define NALICE 7  // number of Alice's attributes
#define NBOB   7  // number of Bob's attributes
#define Nd     6  // number required in common

int Alice[NALICE]={7,6,3,4,12,1,9};    // Alice's attributes
int Bob[NBOB]=    {6,3,4,12,5,10,7};   // Bob's attributes

// Check if person has attribute a
int has_attribute(int num,int *attr,int a)
{
	for (int i=0;i<num;i++)
	{
		if (a==attr[i]) return i;
	}
	return -1;
}

void calculateCoefficients(int n, Big* a, Big* coefficients, Big order) {
    coefficients[0] = 1;

    for (int i = 0; i < n; i++) {
        for (int j = i ; j >= 1; j--) {
            coefficients[j] = modmult(coefficients[j],a[i],order) + coefficients[j - 1];
						coefficients[j] %= order;
        }
        coefficients[0] = modmult(coefficients[0],a[i],order);
        coefficients[i+1] = 1;
    }
}

int main()
{
  int i,j;
  cout<<endl<<"Universe of attributes : 1 - "<<NATTR<<endl<<endl<<"Alice's attributes : ";
  for(i=0;i<NALICE;i++)
    cout<<Alice[i]<<" ";
  cout<<endl<<endl<<"Bob's attributes : ";
  for(i=0;i<NBOB;i++)
    cout<<Bob[i]<<" ";
  cout<<endl<<endl<<"Threshold : "<<Nd<<endl;
  clock_t kgen_start, kgen_finish;
  double kgen_duration;
	PFC pfc(AES_SECURITY);  // initialise pairing-friendly curve
	miracl *mip=get_mip();  // get handle on mip (Miracl Instance Pointer)
	Big order=pfc.order();  // get pairing-friendly group order

	time_t seed;            // crude randomisation
	time(&seed);
    irand((long)seed);

// setup - for 20 attributes 1-20
  //cout<<"setup - for 20 attributes 1-20"<<endl;
  double START,END; 
  START = clock();
  Big tao[2*NATTR];
  for(i=0;i<2*NATTR;i++)
    pfc.random(tao[i]);
  G1 g,u;
  G2 h,paremh[2*NATTR];
  GT v;
  Big alpha,gamma;
  pfc.random(g);
  pfc.random(h);
	pfc.precomp_for_mult(h);  // h is fixed, so precompute on it
  pfc.random(alpha);
  pfc.random(gamma);
  u=pfc.mult(g,modmult(alpha,gamma,order));
  v=pfc.pairing(h,pfc.mult(g,alpha));
  for(i=0;i<2*NATTR;i++)
  {
    paremh[i]=pfc.mult(h,modmult(alpha,pow(gamma,i,order),order));
		pfc.precomp_for_mult(paremh[i],TRUE);// paremh[i] are system params, so precompute on them
     									// Note second parameter indicates that all multipliers
										  // must  be <=2*AES_SECURITY bits, which may be shorter
										  // than the full group size.
	}
  END = clock();
  cout << endl << "Setup：" << (END - START) / CLOCKS_PER_SEC << "s" << endl;
// key extration for Bob
  //cout<<"key extration for Bob"<<endl;
  START = clock();
  Big r;
  pfc.random(r);
  G1 sk1[NBOB];
  G2 sk2[NATTR],sk3;
  for(i=0;i<NBOB;i++)
  {
    sk1[i]=pfc.mult(g,moddiv(r,gamma+tao[Bob[i]],order));
		pfc.precomp_for_mult(sk1[i]);   // Bob precomputes on his private key
  }
  for(i=0;i<NATTR-1;i++)
  {
    sk2[i]=pfc.mult(h,modmult(r,pow(gamma,i,order),order));
		pfc.precomp_for_mult(sk2[i]);   // Bob precomputes on his private key
  }
  sk3=pfc.mult(h,moddiv(r-1,gamma,order));
  END = clock();
  cout << endl << "Key extration for Bob：" << (END - START) / CLOCKS_PER_SEC << "s" << endl;

// encryption by Alice
  //cout<<"encryption by Alice"<<endl;
  START = clock();
  G1 c1;
  G2 c2;
  GT K;
  Big kappa,c3;
  pfc.random(kappa);
  c1=pfc.mult(u,modmult(kappa,(Big)(-1),order));
  Big cof[NATTR+Nd];
  Big sattr[NATTR+Nd-1];
  for(i=0;i<NALICE;i++)
    sattr[i]=tao[Alice[i]];
  for(j=0;j<NATTR+Nd-1-NALICE;j++)
  {
    sattr[i]=tao[j+NATTR+1];
    i++;
  }
  calculateCoefficients(NATTR+Nd-1, sattr, cof, order);
	c2=pfc.mult(paremh[0],cof[0]);
  for(i=1;i<NATTR+Nd;i++)
    c2=c2+pfc.mult(paremh[i],cof[i]);
  c2=pfc.mult(c2,kappa);
  K=pfc.power(v,kappa);//K
  Big M;
	mip->IOBASE=256;
	M=(char *)"test message";
	cout << endl << "Message to be encrypted=   " << M << endl;
	mip->IOBASE=16;
	c3=lxor(M,pfc.hash_to_aes_key(K));
  END = clock();
  cout << endl << "Encryption by Alice：" << (END - START) / CLOCKS_PER_SEC << "s" << endl;

// Decryption by Bob
  //cout<<"Decryption by Bob"<<endl;
  START = clock();
  int k=0,n,S[NBOB],BIndex[NBOB];
  for (j=0;j<NBOB;j++)
	{ // check for common attributes
		i=Bob[j];
		n=has_attribute(NALICE,Alice,i);
		if (n<0) continue;  // Bob doesn't have it
		S[k]=i;   // S is the set of common attributes
    //cout<<"AS "<<k<<" = "<<i<<endl;
    BIndex[k]=j;// index of common attributes in bob
		k++;
	}
	if (k<Nd)
	{
		cout << "Bob does not have enough attributes in common with Alice to decrypt successfully" << endl;
		//exit(0);
	}
  GT L;
  
  G1 p[Nd+1][Nd+1];
  for(i=0;i<Nd;i++)
    for(j=i+1;j<=Nd;j++)
      if(i==0)
        p[i][j]=sk1[BIndex[j-1]];
      else
        p[i][j]=pfc.mult(p[i-1][i]+(-p[i-1][j]),moddiv((Big)1,(tao[S[j-1]]-tao[S[i-1]]),order));
        
  L=pfc.pairing(c2,p[Nd-1][Nd]);
  
	for(i=0;i<NATTR+Nd-1-NALICE;i++)
		sattr[i]=tao[i+NATTR+1];
  for(j=0;j<NALICE;j++)
		if(has_attribute(Nd,S,Alice[j])==-1)
		{
			sattr[i]=tao[Alice[j]];
			i++;
		}
	calculateCoefficients(NATTR-1, sattr, cof, order);
	G2 req1 = pfc.mult(sk2[0],cof[1]);
	for(i=1;i<NATTR-1;i++)
		req1 = req1 + pfc.mult(sk2[i],cof[i+1]);
	GT eq1=pfc.pairing(req1,c1);
  eq1=eq1*L;
	GT eq2=pfc.pairing(sk3,c1);
	K=pfc.power(eq1,moddiv(1,cof[0],order))*eq2;

	M=lxor(c3,pfc.hash_to_aes_key(K));
	mip->IOBASE=256;
	cout << endl << "Decrypted message=         " << M << endl;
  END = clock();
  cout << endl << "Decryption by Bob：" << (END - START) / CLOCKS_PER_SEC << "s" << endl << endl;
  return 0;
}
