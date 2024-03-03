/*
   My ABE - Javier Herranz, Fabien Laguillaumie, and Carla R`afols
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

   In linux: replace "cl /O2 /GX" to "g++"

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
#define Nd     5  // number required in common

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

// Lagrange interpolation
Big lagrange(int i,int *S,int d,Big order)
{
	int j,k;
	Big z=1;
	for (k=0;k<d;k++)
	{
		j=S[k];
		if (j!=i) z=modmult(z,moddiv(order-j,(Big)(i-j),order),order);
	}
	return z;
}

void calculateCoefficients(int n, Big* a, Big* coefficients) {
    coefficients[0] = 1;

    for (int i = 0; i < n; i++) {
        for (int j = i ; j >= 1; j--) {
            coefficients[j] = coefficients[j] * a[i] + coefficients[j - 1];
        }
        coefficients[0] =  coefficients[0] * a[i];
        coefficients[i+1] = 1;
    }
}

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

// setup - for 20 attributes 1-20
  cout<<"setup - for 20 attributes 1-20"<<endl;
  G1 p1,u;
  G2 p2,paremh[2*NATTR];
  GT v;
  int i,j,d[NATTR];
  Big alpha,gamma;
  for(i=0;i<NATTR-1;i++) // default attributes 21-39
    d[i]=i+21;
  pfc.random(p1);
  pfc.random(p2);
  pfc.random(alpha);
  cout<<"alpha = "<<alpha<<endl;
  pfc.random(gamma);
  cout<<"gamma = "<<gamma<<endl;
  Big tao[2*NATTR];
  for(i=0;i<2*NATTR;i++)
    pfc.random(tao[i]);
  u=pfc.mult(p1,alpha);
  u=pfc.mult(u,gamma);
  cout<<"u"<<endl;
  v=pfc.power(pfc.pairing(p2,p1),alpha);
  cout<<"125"<<endl;
  for(i=0;i<2*NATTR;i++)
  {
    paremh[i]=pfc.mult(p2,alpha);
    for(j=0;j<i;j++)
      paremh[i]=pfc.mult(paremh[i],gamma);
  }  
// key extration for Alice
  cout<<"key extration for Alice"<<endl;
  Big ra,temp1;
  pfc.random(ra);
  G1 sk1a[NALICE];
  G2 sk2a[NATTR],sk3a;
  for(i=0;i<NALICE;i++)
  {
    temp1=ra/(gamma+tao[Alice[i]]);
    sk1a[i]=pfc.mult(p1,temp1);
  }
  cout<<"145"<<endl;
  for(i=0;i<NATTR-1;i++)
  {
    sk2a[i]=pfc.mult(p2,ra);
    for(j=0;j<i;j++)
      sk2a[i]=pfc.mult(sk2a[i],gamma);
  }
  cout<<"153"<<endl;
  temp1=(ra-1)/gamma;
  sk3a=pfc.mult(p2,temp1);
 
// key extration for Bob
  cout<<"key extration for Bob"<<endl;
  Big rb;
  pfc.random(rb);
  G1 sk1b[NBOB];
  G2 sk2b[NATTR],sk3b;
  for(i=0;i<NBOB;i++)
  {
    temp1=rb/(gamma+tao[Bob[i]]);
    sk1b[i]=pfc.mult(p1,temp1);
  }
  cout<<"160"<<endl;
  for(i=0;i<NATTR-1;i++)
  {
    sk2b[i]=pfc.mult(p2,rb);
    for(j=0;j<i;j++)
      sk2b[i]=pfc.mult(sk2b[i],gamma);
  }
  cout<<"167"<<endl;
  temp1=(rb-1)/gamma;
  sk3b=pfc.mult(p2,temp1);

// encryption by Alice
  cout<<"encryption by Alice"<<endl;
  G1 c1;
  G2 c2;
  GT K;
  Big kappa,c3;
  pfc.random(kappa);
  c1=pfc.mult(u,-kappa);
  cout<<"181"<<endl;
  Big cof[NATTR+Nd];
  Big sattr[NATTR+Nd-1];
  for(i=0;i<NALICE;i++)
    sattr[i]=tao[Alice[i]];
  for(j=0;j<NATTR+Nd-1-NALICE;j++)
  {
    sattr[i]=tao[d[j]];
    i++;
  }
  calculateCoefficients(NATTR+Nd-1, sattr, cof);
  cout<<"203"<<endl;
  c2=pfc.mult(paremh[0],cof[0]);
  cout<<"204"<<endl;
  for(i=1;i<NATTR+Nd-1;i++)
    c2=c2+pfc.mult(paremh[i],cof[i]);
  cout<<"206"<<endl;
  c2=pfc.mult(c2,kappa);
  /*
  c2=pfc.mult(p2,kappa);
  c2=pfc.mult(c2,alpha);
  c2=pfc.mult(c2,gamma+tao[Alice[0]]);
  for(i=1;i<NALICE;i++)
  {
    cout<<"i="<<i<<endl;
    c2=pfc.mult(c2,gamma+tao[Alice[i]]);
  }
  cout<<"181"<<endl;
  c2=pfc.mult(c2,gamma+tao[d[0]]);
  for(i=1;i<(NATTR+Nd-NALICE-1);i++)
    c2=pfc.mult(c2,gamma+tao[d[i]]);
  */
  cout<<"186"<<endl;
  cout<<"187"<<endl;
  K=pfc.power(v,kappa);
  cout<<"189"<<endl;
  Big M;
	mip->IOBASE=256;
	M=(char *)"test message"; 
	cout << "Message to be encrypted=   " << M << endl;
	mip->IOBASE=16;
	c3=lxor(M,pfc.hash_to_aes_key(K));

// Decryption by Bob
  cout<<"Decryption by Bob"<<endl;
  int k=0,n,S[NBOB],BIndex[NBOB];
  for (j=0;j<NBOB;j++)
	{ // check for common attributes
		i=Bob[j];
		n=has_attribute(NALICE,Alice,i);
		if (n<0) continue;  // Alice doesn't have it
		S[k]=i;   // S is the set of common attributes
    BIndex[k]=j;// index of common attributes in bob
		k++;
	}
	if (k<Nd)
	{
		cout << "Bob does not have enough attributes in common with Alice to decrypt successfully" << endl;
		exit(0);
	}
  GT L;
  G2 TEMP2=pfc.mult(p2,ra);
  Big prod4=1,prod5=1;
  cout<<"223"<<endl;
  
  cout<<"225"<<endl;
  //L= pfc.pairing(c2,pfc.mult(p1,ra/prod3));
  G1 p[Nd][Nd];
  for(i=0;i<Nd;i++)
  {
    for(j=i+1;j<=Nd;j++)
    {
      cout<<"DP i = "<<i<<" , j = "<<j<<endl;
      if(i==0)
        p[i][j]=sk1b[BIndex[j-1]];
      else
      {
        p[i][j]=pfc.mult((p[i-1][i]+pfc.mult(p[i-1][j],-1)),1/(tao[S[j-1]]-tao[S[i-1]]));
        cout<<"DP i = "<<i<<" , j = "<<j<<endl;
      }
    }
  }
  cout<<"243"<<endl;
  L=pfc.pairing(c2,p[Nd-1][Nd]);
  cout<<"226"<<endl;
  int Deset[NATTR];
  int NDeset=0;
  for(j=1;j<NATTR*2;j++)
  {
    if((has_attribute(NALICE,Alice,j)||has_attribute(NATTR+Nd-1-NALICE,d,j))&&!has_attribute(Nd,S,j)) 
    {
      Deset[NDeset]=j;
      cout<<"Dset["<<NDeset<<"] = "<<j<<endl;
      NDeset++;
    }
  }
  Big poly=1;
  for(j=0;j<NDeset;j++)
  {
    prod4*=(gamma+Deset[j]);
    prod5*=Deset[j];
  }
  cout<<"244"<<endl;
  poly=(prod4-prod5)/gamma;
  GT eq1=pfc.pairing(pfc.mult(p2,ra*poly),c1)*L;
  GT eq2=pfc.pairing(pfc.mult(p2,(ra-1)/gamma),c1);
  K= pfc.power(eq1,1/prod5)*eq2;
	M=lxor(c3,pfc.hash_to_aes_key(K));
	mip->IOBASE=256;
	cout << "Decrypted message=         " << M << endl;

    return 0;
}
