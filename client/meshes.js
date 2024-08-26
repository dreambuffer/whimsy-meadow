// meshes.js
import * as THREE from 'three';

export function createMorph(scene, data) {

  const basePath = "client/data/sprites/";
  const file = data["Kaizomorph"]["Name"];
  const fullPath = basePath + file + ".png";

  const morphTexture = new THREE.TextureLoader().load(fullPath);
  const material = new THREE.SpriteMaterial({
    map: morphTexture,
  });
  morphTexture.minFilter = THREE.NearestFilter;
  morphTexture.magFilter = THREE.NearestFilter;
  const morph = new THREE.Sprite(material);
  morph.scale.set(6, 6, 1);
  let pos = new THREE.Vector3(
    data["pos"]["X"],
    data["pos"]["Y"],//          sprite.position.y = Math.sin(1 * clock.getElapsedTime()) * 2;
    data["pos"]["Z"],
  );
  morph.position.copy(pos); // lerp other clients besides player
  return morph
}

export function createGrassPlane(scene) {
  const loader = new THREE.TextureLoader();
  const texture = loader.load("client/data/Png/grass_17.png");
  texture.repeat.set(0.1, 0.1);
  texture.minFilter = THREE.NearestFilter;
  texture.magFilter = THREE.NearestFilter;

  const material = new THREE.MeshBasicMaterial({ map: texture });

  const geometry = new THREE.PlaneGeometry(2000, 2000); // You can adjust the size as needed

  const plane = new THREE.Mesh(geometry, material);
  plane.rotation.x = -Math.PI / 2;
  plane.position.y += 8;

  scene.add(plane);
}
// Define your meshes
export const createSphere = () => {
  const geometry = new THREE.SphereGeometry(1, 32, 32);
  const material = new THREE.MeshBasicMaterial({ color: 0xff0000 });
  return new THREE.Mesh(geometry, material);
};

export const createBox = () => {
  const geometry = new THREE.BoxGeometry(1, 1, 1);
  const material = new THREE.MeshBasicMaterial({ color: 0x00ff00 });
  return new THREE.Mesh(geometry, material);
};


function getRandomColor(start, range) {
  const letters = "0123456789ABCDEF";
  let color = "#";

  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * range) + start];
  }

  return color;
}

// Function to create a flower group
export function createFlower() {
  // Function to create a part of the flower
  const createPart = (
    geometry,
    material,
    position = new THREE.Vector3(),
    rotation = new THREE.Euler(),
  ) => {
    const part = new THREE.Mesh(geometry, material);
    part.position.copy(position);
    part.rotation.copy(rotation);
    return part;
  };

  // Create middle part of the flower
  const middle = createPart(
    new THREE.SphereGeometry(1.5, 32, 32),
    new THREE.MeshBasicMaterial({ color: 0xffd700 }), // Gold color for the middle
    new THREE.Vector3(0, 12, 0),
  );

  // Create petals of the flower
  const petalGeometry = new THREE.CircleGeometry(2.7, 32);
  const randColor = getRandomColor(3, 13);
  const petalMaterial = new THREE.MeshBasicMaterial({ color: randColor }); // Deep Pink color for the petals

  const petals = Array.from({ length: 4 }, (_, index) => {
    const rotationAngle = (Math.PI / 4) * (index * 2);

    const position = new THREE.Vector3(
      4 * Math.cos(rotationAngle),
      12,
      4 * Math.sin(rotationAngle),
    );

    const rotation = new THREE.Euler(-Math.PI / 2, 0, rotationAngle);
    const x = createPart(
      petalGeometry,
      petalMaterial,
      position,
      rotation,
    );

    return x;
  });

  // Create stem of the flower
  const stem = createPart(
    new THREE.CylinderGeometry(0.7, 0.7, 2, 32),
    new THREE.MeshBasicMaterial({ color: 0x228b22 }), // Forest Green color for the stem
    new THREE.Vector3(0, 10, 0),
  );

  // Create the flower group and add parts to it
  const flower = new THREE.Group();
  flower.add(middle, ...petals, stem);

  return flower;
}

// Function to change the color of a petal of the flower
export function changePetalColor(flowerGroup, petalIndex, newColor) {
  // Ensure petalIndex is valid
  if (petalIndex >= 0 && petalIndex < flowerGroup.children.length - 1) {
    const petalMesh = flowerGroup.children[petalIndex + 1]; // Index 0 is the middle, so petal meshes start from index 1
    if (petalMesh.material instanceof THREE.MeshBasicMaterial) {
      petalMesh.material.color.set(newColor);
    }
  }
}

// Function to generate and add flowers to the scene
export function addFlowers(scene) {
  const fieldGroup = new THREE.Group();
  const flowerMap = new Map();

  // Create and add flowers to the field group
  for (let i = 0; i < 200; i++) {
    const flower = createFlower();

    const radius = 950 * Math.sqrt(Math.random()); // Ensures a uniform distribution within a circle
    const angle = Math.random() * 2 * Math.PI;

    const randomX = radius * Math.cos(angle);
    const randomZ = radius * Math.sin(angle);
    const randomY = Math.random() * 0.5;

    flower.position.set(randomX, randomY, randomZ);
    flowerMap.set(flower, flower.position.clone());
    fieldGroup.add(flower);
  }

  // Add the field group to the scene
  scene.add(fieldGroup);
}
//
      export function createTorus(radius, tube, color, position) {
        const geometry = new THREE.TorusGeometry(radius, tube, 16, 100);
        const material = new THREE.MeshPhongMaterial({
          color: color,
          transparent: true,
          opacity: 0.8,
          emissive: 0x0000ff,
        });
        const torus = new THREE.Mesh(geometry, material);
        torus.position.copy(position);
        return torus;
      }



      export function addToruses(scene) {
        const localToruses = [];
      
        for (let i = 0; i < 10; i++) {
          const torus = createTorus(
            16,
            1.5,
            0x2688ff,
            new THREE.Vector3(
              Math.random() * 4 - 2,
              16 + i * 8,
              Math.random() * 4 - 2
            )
          );
          torus.rotation.x = -Math.PI / 2;
          torus.visible = false;
      
          scene.add(torus);
          localToruses.push(torus);
        }
      
        return localToruses;
      }
      

      /*
     const r1 = 150;
      const r2 = 278;
      const thetaLength = Math.PI / 2;
      const geometry = new THREE.RingGeometry(r1, r2, 64, 1, 0, thetaLength);
      const mat = new THREE.MeshBasicMaterial({
        color: 0xffff00,
        side: THREE.DoubleSide,
      });
      const ring = new THREE.Mesh(geometry, mat);
      ring.rotation.x = -Math.PI / 2;
      ring.rotation.z = -Math.PI / 2;
      ring.position.z += 250;
      ring.position.x -= r1 + 64;
      scene.add(ring);
      */