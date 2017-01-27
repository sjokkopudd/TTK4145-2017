# Presentasjon

## Nettverkstopologi
Vi ønsker å implementere et peer to peer nettverk der alle heisene oppfører seg deterministisk på felles delt informasjon. Vi ønsker å gjøre dette ved hjelp av et UDP nettverk. UDP er en lettere protokoll sammenlignet med TCP og vi føler ikke at vi har behov for alt TCP kan tilby. Da UDP ikke har en innebygd beskjedanerkjennelsesfunksjon er dette noe vi vil implementere på applikasjonsnivå.

Alle heisene transmitter en datapakker ved en "event" som inneholder hva som har skjedd og en melding-id.

Når en heis mottar en datapakke oppdaterer den et lokalt kart med som inneholder all informasjon om de andre heisene og knappene. Ved å iterere seg gjennom kartet kan hver heis regne seg ut til hva den burde gjøre.

## Moduldesign
Vi ønsker å dele systemet opp i så uavhengige moduler som mulig. HVer modul skal kunne utføre sin oppgave med minst mulig avhengighet fra de andre modulene.

### Nettverksmodul - udpNettwork.go

Oppgaven til denne modulen er å håndtere kommunikasjon med de andre heisene. Vi ønsker at dette skal skje i en egen rutine som leser en buffer som kan fylles av de andre modulene. Til dette vil følgende interface funksjoner være nødvendig.

* `void startNetworkComunication()` 
Starter kommunikasjonen over nettverket. Dette er "hjernen" i denne modulen. Både sending og reciving vil skje her. 

* `void fillTransmitBuffer(map)`
Laster kartet inn på nettverksmodulens buffer som sendes ut på nettverket av interne funksjoner.


### Kartmodul - elevatorMap.go

Denne modulen har som oppgave å holde oversikt over alle heisene sin posisjon og retning, og samtidig vite hvilke knapper som er trykket inn. Den skal også sørge for en backup-log til eventuell re-start.

* `void newEvent(event)` 
Legger en ny hendelse inn i kartet.

* `map getMap()`
Returnerer en kopi av det nåværende kartet


### Knapper, lys og sensorer - localIO.go

localIO.go har ansvar for å sjekke om noen nye knapper har blitt trykket på, og i så fall si ifra til elevtorMap.go at det har kommet en ny hendelse. Det samme gjelder for heisens posisjon. Når en ny hendelse blir lagt inn i kartet slås de riktige lysene på. 

* `void startSensorPolling()`

### Oppgavebehandler - taskHandler.go

Denne modulen beregner kostfunksjonen basert på kartets nåværende tilstand. Basert på dette gir den instruksjoner til motormodulen. Siden alle heisene har deterministisk oppførsel basert på kostfunksjonen vil taskHandler.go alltid vite hvilke oppdrag heisen skal utføre. 

* `void startElevator()`

### Heisdriver - motor.go

Setter hastighet og retning for heisen. Åpner dørene.

* `void setDirection(direction)`



## Feilhåndtering 

Ved oppstart vil alle heisene dele hele sitt kart med de andre heisene.

* Hva skjer om en upd pakke ikke kommer frem?

En mottatt udp packet vil verifiseres med en ACK sendt i retur. Om ACK ikke kommer frem vil den originale pakken sendes på nytt frem til ACK blir mottatt. Etter ett gitt antall mislykkede ACK vil en heis markeres som død.

* Hva skjer om en heis mister nettverksforbinelsen?

Om en heis mister nettverksforbinelsen vil den ikke motta meldinger fra de andre heisene. Den vil da markere de andre heisene som døde og vice versa. En død heis vil ha uendelig kost og vil derfor være uegnet til å utføre oppgaver sett i de andre heisene sine øyne. Dette betyr at nettverket opererer videre med n-1 heiser og den døde heisen vil fungere som en solo-heis.

* Hva skjer om noen drar ut stikkontakten? 

Alle heiser vil ha lokale kopier av oppgavekartet sitt lagret på hdd'en. Ved oppstart av en heis vil den dele dette kartet med de andre heisene på nettverket og de andre heisene vil sende sitt kart i retur. Dette er en rutine som alltid utføres ved oppstart av en heis. På denne måten vil en ny-oppstartet heis få kartet sitt fylt ut med oppgaver fra de andre heisene.

* Hvordan håndterer vi et "flipped bit"?

Ved å bruke checksum kan vi verifisere om den mottatte pakken er lik den sendte. Om checksum ikke stemmer sendes ikke ACK.