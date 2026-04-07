---
name: blender-test-sicuro
description: Skill per testare l'integrazione con Blender tramite MCP. Genera script Python per Blender ma richiede SEMPRE l'approvazione umana prima dell'esecuzione.
---

# Obiettivo
Generare codice Python per manipolare la scena di Blender in base alle richieste dell'utente, garantendo la massima sicurezza.

# Istruzioni passo-passo
1. Ricevi la richiesta dell'utente su cosa modellare o modificare.
2. Genera lo script Python utilizzando le API `bpy` di Blender.
3. BLOCCO DI SICUREZZA: Non inviare lo script al server MCP.
4. Mostra all'utente l'intero script Python generato in un blocco di codice.
5. Chiedi esplicitamente: "Posso inviare questo script a Blender?".
6. Solo dopo aver ricevuto un "Sì" esplicito dall'utente, esegui il comando nel terminale per inviare lo script a Blender:
7. Usa questo comando: `cat << 'EOF' | python3 .agent/skills/blender-test/scripts/send.py` (sostituendo EOF con il codice python generato).

# Regole di Sicurezza (NON VIOLARE MAI)
* Non importare mai i moduli `os`, `sys`, `subprocess` o `shutil`.
* Non scrivere istruzioni che eliminano file o formattano directory.
* Limita il codice esclusivamente alle manipolazioni della geometria e dei materiali (modulo `bpy`).