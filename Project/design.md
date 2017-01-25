# Design

## Nettverkdesign 

Vi ønsker å implementere et peer to peer nettverk der alle heisene oppfører seg deterministisk på felles delt informasjon.

Alle heisene transmitter datapakker ved en fast frekvens med følgende innhold:

* Oppgavetabell 
* Heisens posisjon, retning og om den lever
* Timestamp

Når en heis mottar en datapakke oppdaterer den et "kart" som ligger lokalt. Kartet oppdateres kontinuerlig med de andre heisenes posisjon og retning. Gjennom å sammenligne kartet og oppgavetabellen kan hver heis regne ut hvem som burde ta hvilket oppdrag. 

## Kostfunksjon og tilstandskart.

For å finne ut hvilken heis som skal gjøre hva vil vi bruke en kostfunksjon. Kostfunksjonen iterer seg gjennom kartet og tilegner en oppgave til heisen med lavest kostnad. Denne tilegningen skjer lokalt; det sendes altså ikke noen beskjed om hvilken heis som gjør hva over nettverket. 

Hvis all informasjon i kartet er riktig vil heisene handle etter samme regler og være enige om hvem som skal gjøre hva. Om en heis ikke er den heisen med lavest kostnad til en oppgave vil den ignorere oppgaven og anta at noen andre tar seg av det.

Om en heis ikke har sendt ut isAlive på en stund vil heisen markeres som død i karet. Dette gir automatisk uendelig kostnad. 

Hvis en heis er underveis med en oppgave når det oppstår en oppgave med lavere kostnad vil heisen utføre den nye oppgaven først. Et eksempel på dette er om heisen er på vei fra 2. etasje til 4. etasje når det trykkes opp i 3. etasje. Det vil da være naturlig å stoppe i 3. 