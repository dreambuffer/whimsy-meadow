
const raceTrackRadius = 50;
const raceTrackWidth = 10;
const raceTrackGeometry = new THREE.TorusGeometry(raceTrackRadius, raceTrackWidth, 64, 100);
const raceTrackMaterial = new THREE.MeshBasicMaterial({
  color: 0x00ff00,
  wireframe: true
});
const raceTrack = new THREE.Mesh(raceTrackGeometry, raceTrackMaterial);
scene.add(raceTrack);


function generatePointsInCircle(radius, numPoints) {
  const points = [];
  const angleIncrement = (2 * Math.PI) / numPoints;

  for (let i = 0; i < numPoints; i++) {
    const angle = i * angleIncrement;
    const x = radius * Math.cos(angle);
    const y = radius * Math.sin(angle);
    points.push({ x, y });
  }

  return points;
}

// Example: generate 100 points within a circle of radius 666
const pointsArray = generatePointsInCircle(666, 100);
console.log(pointsArray);

/*
    // for p in path, if user position is greater than p +- radius of road, then clamp
    function clampToPath(path, position, minX, maxX, minZ, maxZ) {
      // generateBounds by pushing the min and max points into an array of points
      // then clamp them
      for (let p=0;p<path.length; p++){
      MAXX = path[p].x + maxX
      MAXZ = path[p].z + maxZ
      MINX = path[p].x - minX
      MINZ = path[p].z - minZ
      position.x = THREE.MathUtils.clamp(point.x, MINX, MAXX);
      position.z = THREE.MathUtils.clamp(point.z, MINZ, MAXZ);
      }
  }
  */