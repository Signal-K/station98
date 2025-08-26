//
//  EarthGlobeView.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI
import SceneKit

struct EarthGlobeView: UIViewRepresentable {
    var pads: [Pad]

    func makeUIView(context: Context) -> UIView {
        let globeView = GlobeSceneView()
        globeView.sceneView.scene = makeScene(cameraNode: globeView.cameraNode)
        globeView.sceneView.allowsCameraControl = true
        globeView.sceneView.backgroundColor = UIColor.black
        globeView.sceneView.autoenablesDefaultLighting = true
        globeView.setupTapGesture()
        return globeView
    }

    func updateUIView(_ uiView: UIView, context: Context) {
        guard let globeView = uiView as? GlobeSceneView else { return }
        globeView.sceneView.scene = makeScene(cameraNode: globeView.cameraNode)
    }

    private func makeScene(cameraNode: SCNNode? = nil) -> SCNScene {
        let scene = SCNScene()

        // Earth
        let earth = SCNSphere(radius: 0.8)
        let material = SCNMaterial()
        material.diffuse.contents = UIImage(named: "earth.jpg")
        earth.firstMaterial = material

        let earthNode = SCNNode(geometry: earth)
        scene.rootNode.addChildNode(earthNode)

        // Rotation
        let rotation = CABasicAnimation(keyPath: "rotation")
        rotation.fromValue = NSValue(scnVector4: SCNVector4(0, 1, 0, 0))
        rotation.toValue = NSValue(scnVector4: SCNVector4(0, 1, 0, Float.pi * 2))
        rotation.duration = 30
        rotation.repeatCount = .infinity
        earthNode.addAnimation(rotation, forKey: "spin")

        // Pads
        for pad in pads {
            let pin = SCNSphere(radius: 0.01)
            let pinMaterial = SCNMaterial()
            pinMaterial.diffuse.contents = UIColor.systemRed
            pin.firstMaterial = pinMaterial

            let pinNode = SCNNode(geometry: pin)
            let position = latLngToXYZ(lat: pad.latitude, lng: pad.longitude, radius: 0.8)
            pinNode.position = position
            pinNode.name = pad.name
            // Add look-at constraint toward camera if available
            if let camera = cameraNode {
                let constraint = SCNLookAtConstraint(target: camera)
                constraint.isGimbalLockEnabled = true
                pinNode.constraints = [constraint]
            }
            earthNode.addChildNode(pinNode)
        }

        // Camera
        let camNode: SCNNode
        if let cameraNode = cameraNode {
            camNode = cameraNode
        } else {
            camNode = SCNNode()
            camNode.camera = SCNCamera()
            camNode.position = SCNVector3(0, 0, 2.5)
        }
        scene.rootNode.addChildNode(camNode)

        return scene
    }

    private func latLngToXYZ(lat: Double, lng: Double, radius: Double) -> SCNVector3 {
        let latRad = lat * Double.pi / 180
        let lngRad = lng * Double.pi / 180
        let x = radius * cos(latRad) * sin(lngRad)
        let y = radius * sin(latRad)
        let z = radius * cos(latRad) * cos(lngRad)
        return SCNVector3(x, y, z)
    }
}

// Custom UIView subclass wrapping SCNView and handling tap gestures
class GlobeSceneView: UIView {
    let sceneView = SCNView()
    let cameraNode = SCNNode()

    override init(frame: CGRect) {
        super.init(frame: frame)
        setup()
    }

    required init?(coder: NSCoder) {
        super.init(coder: coder)
        setup()
    }

    private func setup() {
        sceneView.frame = bounds
        sceneView.autoresizingMask = [.flexibleWidth, .flexibleHeight]
        addSubview(sceneView)
        // Setup camera node
        cameraNode.camera = SCNCamera()
        cameraNode.position = SCNVector3(0, 0, 2.5)
    }

    func setupTapGesture() {
        let tap = UITapGestureRecognizer(target: self, action: #selector(handleTap(_:)))
        sceneView.addGestureRecognizer(tap)
    }

    @objc private func handleTap(_ gesture: UITapGestureRecognizer) {
        let location = gesture.location(in: sceneView)
        let hits = sceneView.hitTest(location, options: nil)
        if let node = hits.first?.node, let name = node.name {
            print("Tapped node: \(name)")
            // Optionally, show info UI here
        }
    }
}
