
## Research Context and Key Concepts ##

####Agent-Based Model (ABM)

>Agent-based modeling is a rule-based, computational modeling methodology that focuses on rules and interactions among the individual components or the agents of the system.... (In biological models) the goal of this modeling method is to generate (artificial) populations of the system components of interest and simulate their interactions in a virtual world. Agent-based models start with rules for behavior and seek to reconstruct, through computational instantiation of those behavioral rules, the observed patterns of behavior[^AnMiDutta-MoscatoVodovotz2009]. 

The emphasis here is on the *artificiality* of the system. We are **not** simulating real-world biological environments, but rather constructing an abstracted *slice* of "the system components of interest", which serves to model a relation (or dynamic) arising from the interaction of the agents in the model – which agents are patterned on types of fauna in the natural world.

See *[Agent-based model in biology](https://en.wikipedia.org/wiki/Agent-based_model_in_biology)* for more.

####Colour Polymorphism (CP)####

From *[The Evolution of Color Polymorphism: Crypticity, Searching Images, and Apostatic Selection](http://digitalcommons.unl.edu/bioscifacpub/52/)* (Bond AB. 2007):

> Polymorphism... (has been) viewed from the outset as an adaptive response to the behavior of visual predators. (...) The argument has three main premises: (*a*) Visual search is subject to processing limitations, in that it is more difficult to search for two or more dissimilar, cryptic prey types simultaneously than to search for only one. (*b*) Predators consequently tend to focus on the most rewarding prey type and effectively overlook the others. (*c*) The result is frequency-dependent, apostatic selection, a higher relative mortality among more abundant prey types that serves to stabilize prey polymorphism. 
> (...)
	> Polymorphic species can be separated into two broad categories that appear to reflect distinguishable evolutionary strategies. One grouping consists of species that occur in many disparate forms, all of which bear a general resemblance to the same background. The number of morphs can range from perhaps three or four to dozens, as in grouse locusts (*Acrydium arenosum*) (Nabours et al. 1933), land snails (*Cepaea nemoralis*) (Cain & Sheppard 1954), or underwing moths (*Catocala micronympha*) (Barnes & McDunnough 1918)... (the) second group of polymorphic species are associated with heterogeneous environments that are divided into large, disparate substrate patches (Bond & Kamil 2006, Merilaita et al. 2001). In such coarse-grained habitats (Levins 1968), an individual can occupy only one substrate type at a time and, because of the disparity in patch appearance, cannot match all substrates equally well. The result is selection for maximum crypticity on only one of the available substrates... (b)ecause each morph specializes on a particular substrate type, these species may be highly selective in choosing a resting location, actively avoiding nonmatching backgrounds (Edmunds 1976, Owen 1980, Sargent 1981).

####Context: CP in *Isocladus armatus*
The ABM of this project is designed as a companion for the research by Dr. James Dale (PhD):
> We will use *Isocladus armatus* as a model system to evaluate... hypotheses for the evolutionary maintenance of extreme colour polymorphism.

>The function of colouration in *I. armatus* is unknown[^Jansen1971]. However in other species of isopods[^Merilaita2001] [^JormalainenMerilaitaTuomi1995] [^Moreira1974] colouration is typically argued to be a defence against predation: it functions as camouflage[^StevensMerilaita2009]. Here we putatively assign three different forms[^StevensMerilaita2009] of visual camouflage adopted by *I. armatus*: background matching, disruptive coloration and masquerade... the *a priori* expectation is that stabilising selection from predation should fix colouration onto the single most effectively cryptic phenotype. The key question is therefore: **why do multiple camouflage strategies co-occur throughout the entire geographic range[^Jansen1971] [^Jansen1968] of this species?**
Indeed, CP is widely distributed among cryptic prey species[^Bond2007]. In some cases, polymorphism involves cryptic prey which specialise on distinct substrate types that occur within heterogeneous environments divided into large, disparate patches[^Merilaita2001] [^FaralloForstner2012]. However, in many other cases, cryptic prey seem almost unrestricted in their morph diversity[^Bond2007] [^WhiteleyOwenSmith1997]. The evolutionary processes that maintain CP in such species is an enduring and unsolved evolutionary problem[^Bond2007].

The functional hypothesis this ABM is designed to serve is that of **Apostatic Selection**:

>**Apostatic selection** views prey polymorphism as the outcome of the perceptual and cognitive characteristics of their predators. **Because it is more difficult to search simultaneously for multiple cryptic prey types, visual predators are expected to focus on the more common forms and overlook the others**[^Bond2007]. Predator specialization on common morphs results in negatively frequency dependent (apostatic) selection favouring rare morphs within a single population[^Bond2007] [^Dale2006] [^Endler1998]. Indeed, the properties of colour variability expressed by *I. armatus* are consistent with a negatively frequency dependent selective mechanism[^DaleLankReeve2001] [^Dale2006], i.e. high degrees of polymodal variation with multiple component parts that appear to co-vary independently. Although commonly invoked to explain the maintenance of CP in many species... validating the apostatic hypothesis in a natural wild system has proven difficult[^Bond2007] [^GrayMcKinnon2007] [^Dale2006].

####Predator – Prey System####

This ABM should describe a dynamic between populations of two types of agent: **Visual Predators** (VP) and **Colour-Polymorphic Prey** (CPP).
The principal factor will be in the *relation* between the two population sets defined by the colour variability of CPP and the perceptual and cognitive characteristics of VP, resulting in **apostatic selection**.

Generally:
>Predators and their prey evolve together. Over time, prey animals develop adaptations to help them avoid being eaten and predators develop strategies to make them more effective at catching their prey. These strategies and adaptations can take many forms including camouflage, mimicry, defensive mechanisms, agility, speed, behaviors and even tool usage that make their job easier.
In nature a balance tends to exist between the predators and prey within an environment. There are a number of factors that can affect it but part of it is the birth and death rates of the predators and prey species. It is logical to expect the two populations to fluctuate in response to the density of one another.
When the prey species is numerous, the number of predators will increase because there is more food to feed them and a higher population can be supported with available resources. As the number of predators begins to increase, the density of the prey population will decrease in response to increased rates of predation. That results in a decrease in the number of predators as the food resource becomes smaller which in turn decreases the rate of predation, allowing the prey species population to flourish again.

See *[Predator-Prey Relationships](https://explorable.com/predator-prey-relationships)* for more.

[^AnMiDutta-MoscatoVodovotz2009]: <sub>An G., Mi Q., Dutta-Moscato J., Vodovotz Y. (2009). Agent-based models in translational systems biology. *Systems Biology and Medicine*, [1(September/October), 159-171](https://dx.doi.org/10.1002/wsbm.45).</sub>

[^Jansen1971]:<sub> Jansen KP. 1971. Ecological studies on intertidal New Zealand Sphaeromatidae (Isopoda-Flabellifera). *Marine Biology* 11, 262-285.</sub>

[^Merilaita2001]: <sub>Merilaita S. 2001. Habitat heterogeneity, predation and gene flow: Colour polymorphism in the isopod, *Idotea baltica*. *Evolutionary Ecology* 15, 103-116.</sub>

[^JormalainenMerilaitaTuomi1995]: <sub>Jormalainen V, Merilaita S. & Tuomi J. 1995. Differential predation on sexes affects colour polymorphism of the isopod *Idotea baltica* (Pallas). *Biological Journal of the Linnean Society* 55, 45-68.</sub>

[^Moreira1974]: <sub>Moreira PS. 1974. Cryptic protective coloration in *Serolis laevis* and *Serolis polaris* (Isopoda: Flabellifera). *Crustaceana (Leiden)* 27, 1-4.</sub>

[^StevensMerilaita2009]: <sub>Stevens M & Merilaita S. 2009. Animal camouflage: current issues and new perspectives. *Philosophical Transactions of the Royal Society B: Biological Sciences* 364, 423-427.</sub>

[^Jansen1968]: <sub>Jansen, KP. 1968. *A Comparative Study of Intertidal Species of Sphaeromidae (Isopoda Flabellifera)*. PhD Thesis: University of Canterbury, New Zealand.</sub>

[^Bond2007]: <sub>Bond AB. 2007. The Evolution of Color Polymorphism: Crypticity, Searching Images, and Apostatic Selection. *Annual Review of Ecology, Evolution and Systematics* 38, 489-514.</sub>

[^FaralloForstner2012]: <sub>Farallo VR & Forstner MRJ. 2012. Predation and the maintenance of color polymorphism in a habitat specialist squamate. *PLoS ONE* 7, e30316.</sub>

[^WhiteleyOwenSmith1997]: <sub>Whiteley DA, Owen DF & Smith DA. 1997. Massive polymorphism and natural selection in *Donacilla cornea* (Poli, 1791) (Bivalve: Mesodesmatidae). *Biological Journal of the Linnean Society* 62, 475-494.</sub>

[^Dale2006]: <sub>Dale J. 2006. Intraspecific variation in coloration. In *Bird Coloration, Vol 2: Function and Evolution* (eds. Hill GE & McGraw KJ). Harvard University Press: Cambridge, MA. 36-86.</sub>

[^Endler1998]: <sub>Endler JA. 1988. Frequency-dependent predation, crypsis and aposematic coloration. *Philosophical Transactions of the Royal Society of London Series B: Biological Sciences* 319, 505-523.</sub>

[^DaleLankReeve2001]: <sub>Dale J, Lank DB & Reeve HK. Signaling individual identity versus quality: A model and case studies with ruffs, queleas, and house finches. *American Naturalist* 158, 75-86 (2001).</sub>

[^GrayMcKinnon2007]: <sub>Gray SM & McKinnon JS. 2007. Linking color polymorphism maintenance and speciation. *Trends in Ecology & Evolution* 22, 71-79.</sub>