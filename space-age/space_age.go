package space

func Age(a float64,p string) float64 {
  if p == "Earth"{
         return   a/31557600 
  } else if p == "Mercury"{
         return a/(31557600 * 0.2408467)
  } else if p == "Venus"{
         return a/(31557600 * 0.61519726)
  } else if p == "Mars"{
         return a/(31557600 * 1.8808158)
  } else if p == "Jupiter"{
         return a/(31557600 * 11.862615)
  } else if p == "Saturn"{
         return a/(31557600 * 29.447498)
  } else if p == "Uranus"{
         return a/(31557600 * 84.016846)
   } else {
         return a/(31557600 * 164.79132)
  }
}
