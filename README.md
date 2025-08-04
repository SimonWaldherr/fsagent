# FSAgent Configuration Builder

**🚀 Live Demo: [https://simonwaldherr.github.io/fsagent/](https://simonwaldherr.github.io/fsagent/)**

Ein professioneller visueller Konfigurations-Generator für FSAgent mit modernster Web-Technologie.

## ✨ Features v2.0

### 🎨 Progressive Web App (PWA)
- **� App-Installation** - Installierbar auf Desktop und Mobile
- **⚡ Offline-Funktionalität** - Arbeiten ohne Internetverbindung  
- **🔄 Auto-Sync** - Intelligentes Caching und Updates
- **🌙 Dark Mode** - Augenschonende dunkle Oberfläche

### 🚀 Enhanced Productivity
- **💾 Auto-Save** - Automatisches Speichern der Workspace
- **⏪ Undo/Redo** - Vollständige Historie (50 Scschritte)
- **📥 Import/Export** - Workspace und Config-Dateien
- **⌨️ Keyboard Shortcuts** - Professionelle Tastenkürzel
- **✅ Real-time Validation** - Sofortige Konfigurationsprüfung

### 📦 Alle FSAgent Module
- **🤖 OpenAI** - KI-gestützte Dateianalyse mit GPT
- **📢 Notifications** - Slack, Discord, Teams Integration
- **🗄️ Database** - SQLite, MySQL, PostgreSQL Support
- **🔧 File Processing** - Hash, Metadaten, Duplikat-Erkennung  
- **📂 FTP/SFTP** - Sicherer Dateitransfer
- **📧 Mail** - E-Mail Benachrichtigungen
- **🌐 HTTP** - REST API Integration

## 🛠 Lokale Entwicklung

```bash
# Repository klonen
git clone https://github.com/SimonWaldherr/fsagent.git
cd fsagent

# GitHub Pages Version lokal starten
cd docs
python3 -m http.server 8080

# Oder mit Node.js
npx serve .
```

## 📁 Verzeichnisstruktur

```text
docs/
├── index.html      # Hauptseite (GitHub Pages)
├── blocks.js       # Blockly Block-Definitionen
├── generator.js    # JavaScript-Code-Generator
└── README.md       # Diese Datei
```

## 🔧 Technische Details

- **Framework:** Google Blockly + Vanilla JavaScript
- **Hosting:** GitHub Pages (statisch)
- **Browser-Support:** Alle modernen Browser
- **Dependencies:** Nur Blockly CDN (keine Server erforderlich)

## 🤝 Beitragen

1. Fork das Repository
2. Feature Branch erstellen (`git checkout -b feature/amazing-feature`)
3. Änderungen committen (`git commit -m 'Add amazing feature'`)
4. Branch pushen (`git push origin feature/amazing-feature`)
5. Pull Request erstellen

## 📝 Lizenz

Dieses Projekt steht unter der MIT Lizenz - siehe [LICENSE](../LICENSE) Datei für Details.

## 🔗 Links

- **Main Repository:** [https://github.com/SimonWaldherr/fsagent](https://github.com/SimonWaldherr/fsagent)
- **Config Builder:** [https://simonwaldherr.github.io/fsagent/](https://simonwaldherr.github.io/fsagent/)
- **Dokumentation:** [README.md](../README.md)
