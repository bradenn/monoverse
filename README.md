# Monoverse

Simulating the universe, one headache at a time.

While the project started off and technically still functions as a physics simulator, its more of a windowed graphics engine since last it was edited in 2021.

This project was retired and replaced by bradenn/nbodies.

### Hexagonal Cellular Automata
![Screen Shot 2021-05-17 at 21.41.16 PM.png](./docs%2FScreen%20Shot%202021-05-17%20at%2021.41.16%20PM.png)

### Barnes Hut N-Bodies approximation
![Screen Shot 2021-04-27 at 11.17.05 AM.png](docs%2FScreen%20Shot%202021-04-27%20at%2011.17.05%20AM.png)

### Generic Instancing
![img_1.png](./docs/img_1.png)

### Unbounded Hexagonal Cellular Automata
![img.png](img.png)

### Performance Monitoring
![img.png](./docs/img.png)



## Notes for Physics Simulations

### Matter Construction

```text
Quark Construction:
    
Baryon Construction:
    Proton:
        Quarks:
            Up1 <-> Gluon <-> Up2
            Up2 <-> Gluon <-> Down1  
            Down1 <-> Gluon <-> Up1  
    Neutron:
        Quarks:
            Up1 <-> Gluon <-> Down1
            Down1 <-> Gluon <-> Down2  
            Down2 <-> Gluon <-> Up1          
```
