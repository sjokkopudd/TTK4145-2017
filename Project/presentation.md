# Presentasjon

## Nettverkstopologi
Vi ønsker å implementere et peer to peer nettverk der alle heisene oppfører seg deterministisk på felles delt informasjon. Vi ønsker å gjøre dette ved hjelp av et UDP nettverk. UDP er en lettere protokoll sammenlignet med TCP og vi har kommet frem til at vi ikke har behov for alt TCP kan tilby. Siden UDP ikke har en innebygget beskjedanerkjennelsesfunksjon er dette noe vi vil implementere på applikasjonsnivå.

Hver heis har et lokalt kart som inneholder informasjon om alle heisenes posisjon, samt hvilke ordrer som ligger i de forskjellige etasjene. Når et "event" inntreffer (knappetrykk, posisjonsendring) lokalt vil en heis oppdatere sitt kart som den så sender til de andre heisene. De andre heisene sammenligner det tilsendte kartet med sitt egent, og legger til informasjon som den selv ikke har. 

Heisene skal som sagt oppføre seg deterministisk. Ved å iterere seg gjennom det lokale kartet kan hver heis regne seg frem til hva den burde gjøre.

## Moduldesign
Vi ønsker å dele systemet opp i så uavhengige moduler som mulig. Hver modul skal kunne utføre sin oppgave med minst mulig avhengighet av de andre modulene.

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
