---
name: Master-skill
description: Genera automaticamente nuove skill per Antigravity. Si attiva quando l'utente chiede di creare un'automazione, una nuova abilità, o di insegnare all'agente un nuovo flusso di lavoro.
---

# Obiettivo
Il tuo compito è analizzare la richiesta dell'utente per un nuovo flusso di lavoro e generare l'infrastruttura completa di una nuova Skill per Antigravity (cartella e file `SKILL.md`).

# Istruzioni passo-passo

1. **Analisi del Requisito:** Quando l'utente ti chiede di creare una nuova skill, chiedi chiarimenti (se necessari) su:
   * Qual è l'obiettivo esatto.
   * Se ci sono script esterni da eseguire.
   * Come gestire eventuali errori.
2. **Progettazione:** Genera un nome breve e descrittivo per la nuova skill (formato `kebab-case`).
3. **Creazione della Struttura:**
   * Crea la directory `.agent/skills/<nome-della-skill>/`.
   * Se la skill richiede script, crea la sottocartella `.agent/skills/<nome-della-skill>/scripts/`.
4. **Scrittura del file SKILL.md:** Crea il file `.agent/skills/<nome-della-skill>/SKILL.md`. Il file DEVE contenere:
   * Frontmatter YAML con `name` e `description` (scritta in terza persona, molto dettagliata).
   * Un'intestazione `# Obiettivo`.
   * Un'intestazione `# Istruzioni passo-passo` con un elenco numerato chiaro.
   * Un'intestazione `# Albero Decisionale` per la gestione degli errori.
5. **Revisione Umana (CRITICO):** Mostra all'utente il contenuto del file `SKILL.md` generato. Chiedi esplicitamente: "Confermi la creazione di questa skill o vuoi apportare modifiche?".
6. **Finalizzazione:** Solo dopo l'approvazione dell'utente, salva definitivamente i file su disco.

# Regole di formattazione per la nuova skill
* Non inserire mai indicazioni vaghe nella skill generata.
* Sii deterministico: usa verbi imperativi ("Controlla", "Esegui", "Leggi").
* Istruisci la nuova skill a non analizzare mai il codice sorgente di eventuali script esterni, ma solo ad eseguirli.

# Albero Decisionale
* Se l'utente rifiuta la bozza -> Chiedi quali sezioni modificare e rigenera il file.
* Se la creazione della cartella fallisce per permessi -> Avvisa l'utente e suggerisci il comando `chmod` appropriato.