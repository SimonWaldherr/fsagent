// FSAgent Blockly Block Definitionen

// Hauptblock für FSAgent Job
Blockly.Blocks['fsagent_job'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("FSAgent Job");
        this.appendDummyInput()
            .appendField("Ordner:")
            .appendField(new Blockly.FieldTextInput("/mnt/server/input/"), "FOLDER");
        this.appendDummyInput()
            .appendField("Match Pattern:")
            .appendField(new Blockly.FieldTextInput(".*\\.txt$"), "MATCH");
        this.appendDummyInput()
            .appendField("Verbose:")
            .appendField(new Blockly.FieldCheckbox("TRUE"), "VERBOSE");
        this.appendDummyInput()
            .appendField("Debounce:")
            .appendField(new Blockly.FieldCheckbox("TRUE"), "DEBOUNCE");
        this.appendDummyInput()
            .appendField("Nur neue Dateien:")
            .appendField(new Blockly.FieldCheckbox("FALSE"), "ONLYNEW");
        this.appendStatementInput("TRIGGER")
            .setCheck("Trigger")
            .appendField("Trigger:");
        this.appendStatementInput("ACTIONS")
            .setCheck("Action")
            .appendField("Aktionen:");
        this.setColour("#3498db");
        this.setTooltip("Hauptkonfiguration für einen FSAgent Job");
    }
};

// Trigger Blöcke
Blockly.Blocks['trigger_fsevent'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Dateisystem Events");
        this.setPreviousStatement(true, "Trigger");
        this.setColour("#e74c3c");
        this.setTooltip("Reagiert auf Dateisystem-Ereignisse (neue/geänderte Dateien)");
    }
};

Blockly.Blocks['trigger_ticker'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Ticker alle")
            .appendField(new Blockly.FieldNumber(2000, 100), "INTERVAL")
            .appendField("ms");
        this.setPreviousStatement(true, "Trigger");
        this.setColour("#e74c3c");
        this.setTooltip("Prüft in regelmäßigen Abständen auf neue Dateien");
    }
};

Blockly.Blocks['trigger_http'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("HTTP Upload auf Port")
            .appendField(new Blockly.FieldNumber(8080, 1), "PORT");
        this.setPreviousStatement(true, "Trigger");
        this.setColour("#e74c3c");
        this.setTooltip("Ermöglicht Datei-Upload über HTTP");
    }
};

// Action Blöcke
Blockly.Blocks['action_sleep'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Warten")
            .appendField(new Blockly.FieldNumber(1000, 0), "TIME")
            .appendField("ms");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Wartet eine bestimmte Zeit");
    }
};

Blockly.Blocks['action_mail'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("E-Mail senden");
        this.appendDummyInput()
            .appendField("Betreff:")
            .appendField(new Blockly.FieldTextInput("Datei verarbeitet"), "SUBJECT");
        this.appendDummyInput()
            .appendField("Text:")
            .appendField(new Blockly.FieldTextInput("Eine neue Datei wurde verarbeitet"), "BODY");
        this.appendDummyInput()
            .appendField("Von:")
            .appendField(new Blockly.FieldTextInput("noreply@company.com"), "FROM");
        this.appendDummyInput()
            .appendField("An:")
            .appendField(new Blockly.FieldTextInput("admin@company.com"), "TO");
        this.appendDummyInput()
            .appendField("CC:")
            .appendField(new Blockly.FieldTextInput(""), "CC");
        this.appendDummyInput()
            .appendField("BCC:")
            .appendField(new Blockly.FieldTextInput(""), "BCC");
        this.appendDummyInput()
            .appendField("SMTP User:")
            .appendField(new Blockly.FieldTextInput("mailuser"), "USER");
        this.appendDummyInput()
            .appendField("SMTP Pass:")
            .appendField(new Blockly.FieldTextInput("password"), "PASS");
        this.appendDummyInput()
            .appendField("SMTP Server:")
            .appendField(new Blockly.FieldTextInput("smtp.company.com"), "SERVER");
        this.appendDummyInput()
            .appendField("SMTP Port:")
            .appendField(new Blockly.FieldNumber(587, 1), "PORT");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Sendet eine E-Mail mit der Datei als Anhang");
    }
};

Blockly.Blocks['action_move'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Datei verschieben nach:")
            .appendField(new Blockly.FieldTextInput("/mnt/server/archive/$file_%Y%m%d%H%M%S"), "DESTINATION");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Verschiebt die Datei an einen anderen Ort");
    }
};

Blockly.Blocks['action_copy'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Datei kopieren nach:")
            .appendField(new Blockly.FieldTextInput("/mnt/server/backup/$file"), "DESTINATION");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Erstellt eine Kopie der Datei");
    }
};

Blockly.Blocks['action_delete'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Datei löschen");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Löscht die Datei");
    }
};

Blockly.Blocks['action_checksize'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Dateigröße prüfen");
        this.appendDummyInput()
            .appendField("Min. Größe:")
            .appendField(new Blockly.FieldTextInput("1KB"), "MIN_SIZE");
        this.appendDummyInput()
            .appendField("Max. Größe:")
            .appendField(new Blockly.FieldTextInput("10MB"), "MAX_SIZE");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Prüft ob die Dateigröße in einem bestimmten Bereich liegt");
    }
};

Blockly.Blocks['action_http_post'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("HTTP POST Request");
        this.appendDummyInput()
            .appendField("URL:")
            .appendField(new Blockly.FieldTextInput("http://api.company.com/upload"), "URL");
        this.appendDummyInput()
            .appendField("Headers:")
            .appendField(new Blockly.FieldTextInput("Content-Type: application/json"), "HEADERS");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Sendet den Dateiinhalt in einem HTTP POST Request");
    }
};

// OpenAI Action Block
Blockly.Blocks['action_openai'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("OpenAI Analyse");
        this.appendDummyInput()
            .appendField("API Key:")
            .appendField(new Blockly.FieldTextInput(""), "API_KEY");
        this.appendDummyInput()
            .appendField("Modell:")
            .appendField(new Blockly.FieldDropdown([
                ["GPT-3.5 Turbo", "gpt-3.5-turbo"],
                ["GPT-4", "gpt-4"],
                ["GPT-4 Turbo", "gpt-4-turbo"],
                ["GPT-4o", "gpt-4o"]
            ]), "MODEL");
        this.appendDummyInput()
            .appendField("Aktion:")
            .appendField(new Blockly.FieldDropdown([
                ["Analysieren", "analyze"],
                ["Übersetzen", "translate"],
                ["Zusammenfassen", "summarize"],
                ["Klassifizieren", "classify"],
                ["Benutzerdefiniert", "custom"]
            ]), "ACTION");
        this.appendDummyInput()
            .appendField("Prompt:")
            .appendField(new Blockly.FieldTextInput("Analyze this file content"), "PROMPT");
        this.appendDummyInput()
            .appendField("Output Datei:")
            .appendField(new Blockly.FieldTextInput(""), "OUTPUT_FILE");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Analysiert Dateien mit OpenAI GPT");
    }
};

// Notification Action Block
Blockly.Blocks['action_notification'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Benachrichtigung senden");
        this.appendDummyInput()
            .appendField("Webhook URL:")
            .appendField(new Blockly.FieldTextInput(""), "WEBHOOK_URL");
        this.appendDummyInput()
            .appendField("Platform:")
            .appendField(new Blockly.FieldDropdown([
                ["Slack", "slack"],
                ["Discord", "discord"],
                ["Microsoft Teams", "teams"]
            ]), "PLATFORM");
        this.appendDummyInput()
            .appendField("Titel:")
            .appendField(new Blockly.FieldTextInput("FSAgent Notification"), "TITLE");
        this.appendDummyInput()
            .appendField("Nachricht:")
            .appendField(new Blockly.FieldTextInput("File processed: $file"), "MESSAGE");
        this.appendDummyInput()
            .appendField("Farbe:")
            .appendField(new Blockly.FieldDropdown([
                ["Erfolg", "good"],
                ["Warnung", "warning"],
                ["Fehler", "danger"]
            ]), "COLOR");
        this.appendDummyInput()
            .appendField("Dateiinhalt anhängen:")
            .appendField(new Blockly.FieldCheckbox("FALSE"), "INCLUDE_FILE");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Sendet Benachrichtigungen an Slack, Discord oder Teams");
    }
};

// Database Action Block
Blockly.Blocks['action_database'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Datenbank Operation");
        this.appendDummyInput()
            .appendField("Treiber:")
            .appendField(new Blockly.FieldDropdown([
                ["SQLite", "sqlite3"],
                ["MySQL", "mysql"],
                ["PostgreSQL", "postgres"]
            ]), "DRIVER");
        this.appendDummyInput()
            .appendField("Verbindung:")
            .appendField(new Blockly.FieldTextInput("./fsagent.db"), "CONNECTION");
        this.appendDummyInput()
            .appendField("Tabelle:")
            .appendField(new Blockly.FieldTextInput("fsagent_files"), "TABLE");
        this.appendDummyInput()
            .appendField("Aktion:")
            .appendField(new Blockly.FieldDropdown([
                ["Log", "log"],
                ["Insert", "insert"],
                ["Import CSV", "import_csv"],
                ["Import JSON", "import_json"]
            ]), "ACTION");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Speichert oder importiert Daten in/aus Datenbank");
    }
};

// File Processing Action Block
Blockly.Blocks['action_fileprocessing'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("Datei-Verarbeitung");
        this.appendDummyInput()
            .appendField("Aktion:")
            .appendField(new Blockly.FieldDropdown([
                ["Hash erstellen", "hash"],
                ["Thumbnail", "thumbnail"],
                ["Validieren", "validate"],
                ["Metadaten extrahieren", "extract_metadata"],
                ["Regex extrahieren", "regex_extract"],
                ["Duplikate finden", "find_duplicates"]
            ]), "ACTION");
        this.appendDummyInput()
            .appendField("Hash Typ:")
            .appendField(new Blockly.FieldDropdown([
                ["MD5", "md5"],
                ["SHA256", "sha256"]
            ]), "HASH_TYPE");
        this.appendDummyInput()
            .appendField("Regex Pattern:")
            .appendField(new Blockly.FieldTextInput(""), "REGEX_PATTERN");
        this.appendDummyInput()
            .appendField("Output Datei:")
            .appendField(new Blockly.FieldTextInput(""), "OUTPUT_FILE");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Erweiterte Datei-Verarbeitungsoperationen");
    }
};

// FTP Action Block
Blockly.Blocks['action_ftp'] = {
    init: function() {
        this.appendDummyInput()
            .appendField("FTP/SFTP Upload");
        this.appendDummyInput()
            .appendField("Host:")
            .appendField(new Blockly.FieldTextInput("ftp.example.com"), "HOST");
        this.appendDummyInput()
            .appendField("Port:")
            .appendField(new Blockly.FieldNumber(22, 1), "PORT");
        this.appendDummyInput()
            .appendField("Protokoll:")
            .appendField(new Blockly.FieldDropdown([
                ["SFTP", "sftp"],
                ["FTP", "ftp"]
            ]), "PROTOCOL");
        this.appendDummyInput()
            .appendField("Benutzername:")
            .appendField(new Blockly.FieldTextInput(""), "USERNAME");
        this.appendDummyInput()
            .appendField("Passwort:")
            .appendField(new Blockly.FieldTextInput(""), "PASSWORD");
        this.appendDummyInput()
            .appendField("Remote Verzeichnis:")
            .appendField(new Blockly.FieldTextInput("/upload/"), "REMOTE_DIR");
        this.appendDummyInput()
            .appendField("Umbenennen:")
            .appendField(new Blockly.FieldTextInput(""), "RENAME");
        this.appendDummyInput()
            .appendField("Lokale Datei behalten:")
            .appendField(new Blockly.FieldCheckbox("TRUE"), "KEEP_FILE");
        this.appendStatementInput("ON_SUCCESS")
            .setCheck("Action")
            .appendField("Bei Erfolg:");
        this.appendStatementInput("ON_FAILURE")
            .setCheck("Action")
            .appendField("Bei Fehler:");
        this.setPreviousStatement(true, "Action");
        this.setNextStatement(true, "Action");
        this.setColour("#27ae60");
        this.setTooltip("Lädt Dateien via FTP oder SFTP hoch");
    }
};
