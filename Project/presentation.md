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

* `void bufferEvent(event)`
Laster innhold inn på modulens buffer som sendes ut på nettverket av interne funksjoner.


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

Hva skjer om en upd pakke ikke kommer frem?

Hva skjer om en heis mister nettverksforbinelsen?

Hva skjer om noen drar ut stikk kontakten? 