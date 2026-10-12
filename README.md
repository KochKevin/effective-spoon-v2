Rein deutsches README

- Verlinkungen
    - Englisches Readme
    - Technische Dokumentationen (Englisch)

# Kiosksystem (effective-spoon)

Ein Kiosk System primär für Getränke welches für Vereine wie Feuerwehren optimiert ist

# TODO: Screenshots und gifs


### Features
- Verwaltung von Getränke und Nutzern
- Erstellen von Ereignissen ("Runde schmeißen")
- Stripe Implemntation zum aufladen des Nutzer Kontos
- Anbindung von RFID Scanner für Identfizierung eines Nutzers und Barcode Lesers für erkennung von Produkten



### Using or Installation

Das ganze System ist in einer Binary enthalten. Bevor diese gestartet werden kann muss ein .env definiert werden (Link zu example.env).
In dieser env wird unteranderem der Stripe Key hinterlegt.
Ansonsten ist das System portabel. Eine Sqlite DB file wird beim systemstart erstellt.
Nutze die prebuilt bins



### Installation on PI

> [!CAUTION]
> Diese Installation ist nicht getestet
> Ich hatte bisher noch nicht die Hardwear um das System so zu installieren
> Ich konnte den groben Stack auf meinem Gerät testen, aber es nicht auf einem PI testen.

> This is designed to use as little of memory as possible.
> Because Hardwear isn't cheap
> An soft goal of this Project is to run it on Raspberry PI Zero 2 W (512mb RAM)

Softwear Stack:
- Raspberry PI OS Lite (Doesn't include an graphical user interface for the OS, to save on memory)
- [Cage](https://github.com/cage-kiosk/cage) um den Browser für die SPA anzuzeigen, kann mit Cage ein einzelnes leichtgewichtiges Fenster erstellt werden
- [Cog](https://github.com/igalia/cog) der eigentliche Browser, basiert auf WebKit
- effective-spoon binary


TODO: Rephrase better
1. [Flash your PI with the RaspberryPI OS Lite](https://www.raspberrypi.com/documentation/computers/getting-started.html#install) and setup SSH
2. Install Cage
3. Install Cog
4. Download the latest effective-spoon binary (look for the latest tag)
5. Setup env file for effective-spoon
6. Setup startup file for pi
    - Let Cage (fullscreen), Cog(fullscreen) and effective spoon start
7. Restart PI to test if everything works




## Architektur

Es handelt sich um eine Klassiches Web Applikation:
Das Frontend ist eine einfache [Nuxt](https://nuxt.com/) [SPA](https://en.wikipedia.org/wiki/Single-page_application)
Das Backend ist in Golang verfasst. Es stellt die API bereit sowie served es die SPA
Als Datenbank wird [SQLite](https://www.sqlite.org/) verwendet, womit das ganze System zu einer ausführbaren Datei zusammenschrumpft. Die verwendet Lib ist [modernc](https://pkg.go.dev/modernc.org/sqlite) variante. Welche nicht auf C abhängigkeiten basiert

Zu Kommuniktion zwischen Frontend und Backend wird eine RESTfull API verwendet. Diese stellt auch einen [SSE](https://html.spec.whatwg.org/multipage/server-sent-events.html) Endpoint bereit um Informationen in das Frontend zu pushen.
Dokumentiert ist die API mit einer OpanAPI Spezifikation



Ein soft Ziel von diesem Projekt war es den memory footprint so klein wie möglich zu halten damit der gesamte Stack theoretisch auf dem Raspberry PI Zero 2 W betrieben werden kann (512 mb)
    - Angesicht der aktuellen hohen RAM und allgemeinen Hardwear Preise (2026)

Dennoch würde ich einen leistungsfähigeren PI empfehlen. Da sich die Hardwear anforderungen in der Zunft ändern können


