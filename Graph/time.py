import csv
import sys
import shlex
import os
import glob
import numpy as np
import matplotlib.pyplot as plt
import shutil
import json
import pylab

fig3 = plt.figure(3, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,1.2])
condition = ['10','20','30','40','50']
setup = [0.26835,0.4816,0.6778,0.8916,1.0892]
key = [0.12833,0.2332,0.3349,0.4388,0.5396]
enc = [0.01078,0.0155,0.0205,0.0257,0.0305]
dec = [0.03009,0.03478,0.0397,0.0448,0.0494]
x = np.arange(len(condition))
width = 0.4
plt.plot(x, setup, 'ro-', linewidth=2, markersize=6, color='red', label='Setup')
plt.plot(x, key, 'ro--', linewidth=2, markersize=6, color='purple', label='KeyGen')
plt.plot(x, enc, 'ro:', linewidth=2, markersize=6, color='green', label='Encrypt')
plt.plot(x, dec, 'ro-.', linewidth=2, markersize=6, color='blue', label='Decrypt')
plt.xticks( x , condition)
plt.xlabel('Number of universe attributes')
plt.ylabel('Average delay(sec)')
#plt.title('Delay between different number of universe attributes')
plt.legend( bbox_to_anchor=(0.2,1))

fig3.show()

fig4 = plt.figure(4, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,0.7])
condition = ['5','10','15','20','25']
setup = [0.640926,0.639958,0.646994,0.64031,0.649864]
key = [0.326198,0.350054,0.37397,0.39947,0.424907]
enc = [0.02059,0.020387,0.02025,0.020574,0.020428]
dec = [0.039759,0.039628,0.039514,0.039878,0.039626]
x = np.arange(len(condition))
width = 0.4
plt.plot(x, setup, 'ro-', linewidth=2, markersize=6, color='red', label='Setup')
plt.plot(x, key, 'ro--', linewidth=2, markersize=6, color='purple', label='KeyGen')
plt.plot(x, enc, 'ro:', linewidth=2, markersize=6, color='green', label='Encrypt')
plt.plot(x, dec, 'ro-.', linewidth=2, markersize=6, color='blue', label='Decrypt')
plt.xticks( x , condition)
plt.xlabel('Number of given user attribute')
plt.ylabel('Average delay(sec)')
#plt.title('Computation cost between different number of user attributes')
plt.legend( bbox_to_anchor=(1,0.87))

fig4.show()

fig5 = plt.figure(5, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,0.7])
condition = ['2','4','6','8','10']
setup = [0.654938,0.653874,0.643219,0.642337,0.653904]
key = [0.399732,0.401498,0.399951,0.399123,0.403753]
enc = [0.018922,0.019955,0.020925,0.021845,0.023098]
dec = [0.035353,0.037942,0.042233,0.048092,0.056494]
x = np.arange(len(condition))
width = 0.4
plt.plot(x, setup, 'ro-', linewidth=2, markersize=6, color='red', label='Setup')
plt.plot(x, key, 'ro--', linewidth=2, markersize=6, color='purple', label='KeyGen')
plt.plot(x, enc, 'ro:', linewidth=2, markersize=6, color='green', label='Encrypt')
plt.plot(x, dec, 'ro-.', linewidth=2, markersize=6, color='blue', label='Decrypt')
plt.xticks( x , condition)
plt.xlabel('Threshold')
plt.ylabel('Average delay(sec)')
#plt.title('Computation cost between different threshold')
plt.legend( bbox_to_anchor=(1,0.87))

fig5.show()

'''

fig = plt.figure(1, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,0.17])
condition = ['Setup','KeyGen','Encrypt','Decrypt']
ecc = [0.00186,0.01348,0.00221,0.00016]
rsa = [0.000316,1.751,0.00052,0.01487]
our = [0.159695,0.056477,0.007012,0.023516]
srsp = [0.0612,0.0355,0.0926,0.0017]
x = np.arange(len(condition))
width = 0.2
plt.bar(x-width, ecc, width, color='purple', label='ECC')
plt.bar(x, rsa, width, color='green', label='RSA')
plt.bar(x+width, srsp, width, color='gray', label='SRSP')
plt.bar(x+width*2, our, width, color='brown', label='Our Scheme')
plt.xticks( x + width / 2, condition)
plt.ylabel('Average delay(sec)')
#plt.title('Delays at various stages between different encryption')
plt.text(1,0.16,'1.751', ha='center', va= 'bottom',fontsize=11)
plt.legend( loc='upper right')

fig1.show()


ecce = 2.37/1000
eccd = 0.16/1000
rsae = 0.52/1000
rsad = 14.87/1000
ourd = 7.01/1000
oure = 23.5/1000
upload = 2.3/1000
download = 1.9/1000
SRSPT = 19.0/1000

fig3 = plt.figure(3, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,0.75])
a = 1
b = 5
c = 10
d = 20
e = 30
condition = [a,b,c,d,e]
ecc = [(ecce+eccd+upload+download)*a,(ecce+eccd+upload+download)*b,(ecce+eccd+upload+download)*c,(ecce+eccd+upload+download)*d,(ecce+eccd+upload+download)*e]
rsa = [(rsae+rsad+upload+download)*a,(rsae+rsad+upload+download)*b,(rsae+rsad+upload+download)*c,(rsae+rsad+upload+download)*d,(rsae+rsad+upload+download)*e]
srsp = [SRSPT*a,SRSPT*b,SRSPT*c,SRSPT*d,SRSPT*e]
our = [(upload+download)*a+ourd+oure,(upload+download)*b+ourd+oure,(upload+download)*c+ourd+oure,(upload+download)*d+ourd+oure,(upload+download)*e+ourd+oure]
x = np.arange(len(condition))
width = 0.4
plt.plot(x, ecc, 'ro-.', linewidth=2, markersize=6, color='red', label='ECC')
plt.plot(x, rsa, 'ro--', linewidth=2, markersize=6, color='purple', label='RSA')
plt.plot(x, srsp, 'ro-', linewidth=2, markersize=6, color='brown', label='SRSP')
plt.plot(x, our, 'ro:', linewidth=2, markersize=6, color='green', label='Our Scheme')
plt.xticks( x, condition)
plt.xlabel('Number of request vehicles')
plt.ylabel('Total delay(sec)')
#plt.title('Total delay between different number of request vehicles')
plt.legend( loc='upper right')

fig3.show()
'''



fig6 = plt.figure(6, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,0.022])
condition = ['PIDGen','Sign','Verify','Trace']
time = [0.005612,0.002826,0.019575,0.001199]
x = np.arange(len(condition))
width = 0.6
plt.bar(x+width/2, time, width, color='blue')
plt.xticks( x + width / 2, condition)
plt.ylabel('Average delay(sec)')
#plt.title('Computation cost between different operation')

fig6.show()

bsfpsi = 0.020186
srspsi = 0.10612
oursi =  0.159695
bsfpwr = 0.01348
srspwr = 0.00455
ourwr = 0.0654
fig7 = plt.figure(7, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,0.17])
condition = ['System initilization','Witness request','Witness response','Verification']
bsfp = [bsfpsi,bsfpwr,0.01021,0.02216]
our = [oursi,ourwr,0.012012,0.019516]
srsp = [srspsi,srspwr,0.006926,0.025017]
x = np.arange(len(condition))
width = 0.3
plt.bar(x-width/2, our, width, color='blue', label='Our Scheme')
plt.bar(x+width/2, bsfp, width, color='green', label='BSFP')
plt.bar(x+width*3/2, srsp, width, color='purple', label='SRSP')
plt.xticks( x + width / 2, condition)
plt.ylabel('Average delay(sec)')
#plt.title('Delays at various stages between different scheme')
plt.legend( loc='upper right')

fig7.show()

trans = 0.04
fig8 = plt.figure(8, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,4.5])
a = 1
b = 5
c = 10
d = 15
e = 20
condition = [a,b,c,d,e]
SRSPe = 4.5/1000
SRSPd = 4.5/1000
BSFPe = 20/1000
BSFPd = 20/1000
ourenc= 60/1000
ourdec= 20/1000
our = [(oursi+ourenc+ourdec+ourwr+trans*4),(oursi+ourenc+ourdec+ourwr+trans*4),(oursi+ourenc+ourdec+ourwr+trans*4),(oursi+ourenc+ourdec+ourwr+trans*4),(oursi+ourenc+ourdec+ourwr+trans*4)]
bsfp = [(BSFPe+BSFPd+3*trans)*a+bsfpsi,(BSFPe+BSFPd+trans*3)*b+bsfpsi,(BSFPe+BSFPd+trans*3)*c+bsfpsi,(BSFPe+BSFPd+trans*3)*d+bsfpsi,(BSFPe+BSFPd+trans*3)*e+bsfpsi]
srsp = [srspsi+(SRSPe+SRSPd+trans*2)*a,srspsi+(SRSPe+SRSPd+trans)*b,srspsi+(SRSPe+SRSPd+trans)*c,srspsi+(SRSPe+SRSPd+trans)*d,srspsi+(SRSPe+SRSPd+trans)*e]
x = np.arange(len(condition))
width = 0.4
plt.plot(x, our, 'ro-', linewidth=2, markersize=6, color='blue', label='Our Scheme')
plt.plot(x, bsfp, 'ro--', linewidth=2, markersize=6, color='green', label='BSFP')
plt.plot(x, srsp, 'ro:', linewidth=2, markersize=6, color='purple', label='SRSP')
plt.xticks( x, condition)
plt.xlabel('Number of request vehicles')
plt.ylabel('Total delay(sec)')
#plt.title('Total delay between different number of request vehicles')
plt.legend( loc='upper right')

fig8.show()

fig9 = plt.figure(9, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,0.55])
for i in range(5):
    our[i] = our[i]/condition[i]
    bsfp[i] = bsfp[i]/condition[i]
    srsp[i] = srsp[i]/condition[i]
x = np.arange(len(condition))
width = 0.4
plt.plot(x, our, 'ro-', linewidth=2, markersize=6, color='blue', label='Our Scheme')
plt.plot(x, bsfp, 'ro--', linewidth=2, markersize=6, color='green', label='BSFP')
plt.plot(x, srsp, 'ro:', linewidth=2, markersize=6, color='purple', label='SRSP')
plt.xticks( x, condition)
plt.xlabel('Number of request vehicles')
plt.ylabel('Average delay(sec)')
#plt.title('Average delay per vehicle between different number of request vehicles')
plt.legend( loc='upper right')

fig9.show()

fig10 = plt.figure(10, figsize=(8,6))
axes = plt.gca()
axes.set_ylim([0,200000])
condition = [1,2,3,4,5]
add = [68741,51641,51641,51641,51641]
ver = [136095,136095,136095,136095,136095]
x = np.arange(len(condition))
width = 0.4
plt.plot(x, add, 'ro-.', linewidth=2, markersize=6, color='red', label='Add message')
plt.plot(x, ver, 'ro--', linewidth=2, markersize=6, color='purple', label='Verify')
plt.xticks( x, condition)
plt.xlabel('Number of additions/verifications (each request)')
plt.ylabel('Gas cost (gas)')
#plt.title('Gas consumption on the smart contract')
plt.legend( loc='upper right')

fig10.show()

input()