# Design

## Nettverkdesign 

Vi ønsker å implementere et peer to peer nettverk der alle heisene oppfører seg deterministisk på felles delt informasjon. Vi ønsker å gjøre dette ved hjelp av et UDP nettverk. UDP er en lettere protokoll sammenlignet med TCP og vi føler ikke at vi har behov for alt TCP kan tilby. Da UDP ikke har en innebygd beskjedanerkjennelsesfunksjon er dette noe vi vil implementere på applikasjonsnivå.

Alle heisene transmitter en datapakker ved en "event" som inneholder hva som har skjedd og en melding-id.

Når en heis mottar en datapakke oppdaterer den et lokalt kart med som inneholder all informasjon om de andre heisene og knappene. Ved å iterere seg gjennom kartet kan hver heis regne seg ut til hva den burde gjøre.

## Kostfunksjon og tilstandskart.

For å finne ut hvilken heis som skal gjøre hva vil vi bruke en kostfunksjon. Kostfunksjonen iterer seg gjennom kartet og tilegner en oppgave til heisen med lavest kostnad. Denne tilegningen skjer lokalt; det sendes altså ikke noen beskjed om hvilken heis som gjør hva over nettverket. 

Hvis all informasjon i kartet er riktig vil heisene handle etter samme regler og være enige om hvem som skal gjøre hva. Om en heis ikke er den heisen med lavest kostnad til en oppgave vil den ignorere oppgaven og anta at noen andre tar seg av det.

Om en heis ikke har sendt ut isAlive på en stund vil heisen markeres som død i karet. Dette gir automatisk uendelig kostnad. 

Hvis en heis er underveis med en oppgave når det oppstår en oppgave med lavere kostnad vil heisen utføre den nye oppgaven først. Et eksempel på dette er om heisen er på vei fra 2. etasje til 4. etasje når det trykkes opp i 3. etasje. Det vil da være naturlig å stoppe i 3. 

## Hva hvis?

* En pakke ikke kommer frem?
 * Heisene vil alltid vente på å få acknowledge signal fra alle de andre heisene som koresponderer til en melding heisen selv har sendt. Dersom heisen mangler acknowledge fra en eller flere heiser etter en gitt tidsperiode forsøker den å sende meldingen på nytt til de mangledne heisene et fastsatt antall ganger. 
  

* Når er en heis død?
 * En heis som ikke sender acknowledge på en melding vil bli ansett som død av heisen som venter på acknowledge. Den ventende heisen vil da si ifra til de andre heisene at den tror at den ikkeresponderende heisen er død. De andre heisene vil forsøke én gang å få tak i den potensielt døde heisen. Dersom de får svar av den vil de be den ventende heisen om å forsøke å få tak i den døde igjen. Dersom dette ikke går vil heisen bli erklært helt død, og kostnaden dens vil bli satt til uendelig. 
  

* En heis faller ut av nettverket?
 * En heis som er erklært død må ikke bare hindres fra å bli tilegnet nye oppdrag, men dens gamle oppdrag må også tas over av de andre heisene. Siden kosten til den døde heisen er satt til uendelig vil dette bare bety at beregningen av oppdragskost vil føre til at en annen heis har lavere kost en den døde, og den vil dermed ta oppdraget.
 * En heis som ikke får kontakt med noen andre heiser vil forsette sin drift alene. Dette vil si at den vil forsøke å utføre alle oppdrag som finnes på kortet, og ta på seg oppdrag som gis fra trykk på dens eget etasjepanel. Med gjevne mellomrom vil den forsøke å resette seg for å få kontakt med nettverket igjen. 
 

* En heis klarer å koble seg på nettverket igjen?
  * Hvis en heis har vært koblet fra nettverket og klarer å koble seg på igjen vil de andre heisene på nettverket få en melding om at heisne er i live. De vil sende acknowledge til den gjenoppståtte heisen, slik at den vet at den er tilbake i drift. Den vil motta kartet fra de andre heisene og legge til sin egen informasjon. Vi er da tilbake i normal operasjon.
  
  
