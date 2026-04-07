{\rtf1\ansi\ansicpg1252\cocoartf2868
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
{\colortbl;\red255\green255\blue255;}
{\*\expandedcolortbl;;}
\paperw11900\paperh16840\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\pard\tx566\tx1133\tx1700\tx2267\tx2834\tx3401\tx3968\tx4535\tx5102\tx5669\tx6236\tx6803\pardirnatural\partightenfactor0

\f0\fs24 \cf0 import sys\
import urllib.request\
import json\
\
# Legge lo script Python generato da Antigravity tramite lo standard input\
script_content = sys.stdin.read()\
\
# Prepara il pacchetto da inviare a Blender\
data = json.dumps(\{'script': script_content\}).encode('utf-8')\
req = urllib.request.Request('http://127.0.0.1:8080', data=data, headers=\{'Content-Type': 'application/json'\})\
\
try:\
    with urllib.request.urlopen(req) as response:\
        print("\uc0\u9989  Invio riuscito:", response.read().decode('utf-8'))\
except urllib.error.URLError as e:\
    print(f"\uc0\u10060  Errore di connessione: Blender \'e8 aperto e lo script \'e8 in esecuzione? Dettagli: \{e\}")}