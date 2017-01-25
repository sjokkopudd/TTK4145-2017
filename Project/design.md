# Design

## Nettverkdesign 

Vi ønsker å implementere et peer to peer nettverk der alle heisene oppfører seg deterministisk på felles delt informasjon.

Alle heisene transmitter datapakker ved en endring (knappetrykk, mangel på acknowledge, status) med følgende innhold:

* Oppgavetabell 
* Heisens posisjon, retning og om den lever
* Meldings-ID

Når en heis mottar en datapakke oppdaterer den et "kart" som ligger lokalt. Kartet oppdateres kontinuerlig med de andre heisenes posisjon og retning. Gjennom å sammenligne kartet og oppgavetabellen kan hver heis regne ut hvem som burde ta hvilket oppdrag. 

## Kostfunksjon og tilstandskart.

For å finne ut hvilken heis som skal gjøre hva vil vi bruke en kostfunksjon. Kostfunksjonen iterer seg gjennom kartet og tilegner en oppgave til heisen med lavest kostnad. Denne tilegningen skjer lokalt; det sendes altså ikke noen beskjed om hvilken heis som gjør hva over nettverket. 

Hvis all informasjon i kartet er riktig vil heisene handle etter samme regler og være enige om hvem som skal gjøre hva. Om en heis ikke er den heisen med lavest kostnad til en oppgave vil den ignorere oppgaven og anta at noen andre tar seg av det.

Om en heis ikke har sendt ut isAlive på en stund vil heisen markeres som død i karet. Dette gir automatisk uendelig kostnad. 

Hvis en heis er underveis med en oppgave når det oppstår en oppgave med lavere kostnad vil heisen utføre den nye oppgaven først. Et eksempel på dette er om heisen er på vei fra 2. etasje til 4. etasje når det trykkes opp i 3. etasje. Det vil da være naturlig å stoppe i 3. 

## Hva hvis?

* En pakke ikke kommer frem?
  * Heisene vil alltid vente på å få acknowledge signal fra alle de andre heisene som koresponderer til en melding heisen selv har sendt. Dersom heisen mangler acknowledge fra en eller flere heiser etter en gitt tidsperiode forsøker den å sende meldingen på nytt til de mangledne heisene et fastsatt antall ganger. 
  
* Når er en heis død?
  * En heis som ikke sender acknowledge på en melding vil bli ansett som død av heisen som venter på acknowledge. Den ventende heisen vil da si ifra til den andre heisene at den tror at den andre heisen er død. De andre heisene vil forsøke én gang å få tak i den potensielt døde heisen. Dersom de får svar av den vil de be den ventende heisen om å forsøke å få tak i den døde igjen. Dersom dette ikke går vil heisen bli erklært helt død, og kostnaden dens vil bli satt til uendelig. 
  
* En heis faller ut av nettverket?
  * 

