# Presentasjon

## Nettverkstopologi
Vi ønsker å implementere et peer to peer nettverk der alle heisene oppfører seg deterministisk på felles delt informasjon. Vi ønsker å gjøre dette ved hjelp av et UDP nettverk. UDP er en lettere protokoll sammenlignet med TCP og vi har kommet frem til at vi ikke har behov for alt TCP kan tilby. Siden UDP ikke har en innebygget beskjedanerkjennelsesfunksjon er dette noe vi vil implementere på applikasjonsnivå.

Hver heis har et lokalt kart som inneholder informasjon om alle heisenes posisjon, samt hvilke ordrer som ligger i de forskjellige etasjene. Når et "event" inntreffer (knappetrykk, posisjonsendring) lokalt vil en heis oppdatere sitt kart som den så sender til de andre heisene. De andre heisene sammenligner det tilsendte kartet med sitt eget, og legger til informasjon som den ikke har selv. 

Heisene skal som sagt oppføre seg deterministisk. Ved å iterere seg gjennom det lokale kartet kan hver heis regne seg frem til hva den burde gjøre.

## Moduldesign
Vi ønsker å dele systemet opp i så uavhengige moduler som mulig. Hver modul skal kunne utføre sin oppgave med minst mulig avhengighet av de andre modulene.

### Nettverksmodul - udpNettwork.go

Oppgaven til denne modulen er å håndtere kommunikasjon med de andre heisene. Vi ønsker at dette skal skje i en egen rutine som leser en channel som kan fylles av elevatorMap.go. Når en ny datapakke blir tatt imot fra andre heiser legger denne modulen pakken inn i en channel som elevatorMap.go kan lese. Til dette vil følgende interface funksjoner være nødvendige:

* `void startNetworkComunication()` 
Starter kommunikasjonen over nettverket. Dette er "hjernen" i denne modulen. Både sending og reciving vil skje her. 


### Kartmodul - elevatorMap.go

Denne modulen har som oppgave å holde oversikt over alle heisene sin posisjon og retning, og samtidig vite hvilke knapper som er trykket inn. Ved en endring i kartet legger den det nye kartet inn i en channel som udpNetwork.go kan lese. Den skal også sørge for en backup-log til eventuell re-start.

* `void newEvent(event)` 
Legger en ny hendelse inn i kartet.

* `map getMap()`
Returnerer en kopi av det nåværende kartet.


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