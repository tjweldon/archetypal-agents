const flock = [];

let alignSlider, cohesionSlider, separationSlider;

const setup = (s) => () => {
    s.createCanvas(800, 400)
    s.background(0)
    separationSlider = s.createSlider(0, 5, 0.1);
    cohesionSlider = s.createSlider(0, 5, 0.1);
    alignSlider = s.createSlider(0, 5, 0.1);
    for (let i = 0; i < 100; i++) {
        flock.push(new Boid(s))
    }
}

let i = 0;
const draw = (s) => () => {
    s.background(0)
    for (let boid of flock) {
        boid.flock(s, flock)
        boid.update(s)
        boid.edges(s)
        boid.show(s)
    }
}


let sketch = (s) => {
    s.setup = setup(s)
    s.draw = draw(s)
}

const sketchInstance = new p5(sketch);