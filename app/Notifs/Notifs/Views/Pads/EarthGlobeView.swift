//
//  EarthGlobeView.swift
//  Station98
//
//  Created by Liam Arbuckle on 26/8/2025.
//

import SwiftUI
import SceneKit

import SwiftUI
import SceneKit

struct EarthGlobeView: UIViewRepresentable {
    var pads: [Pad]
    @Binding var isRotating: Bool
    @Binding var showLabels: Bool
    @Binding var selectedCountries: Set<String>

    func makeUIView(context: Context) -> UIView {
        let globeView = GlobeSceneView()
        
        let camera = SCNCamera()
        camera.zFar = 1000
        camera.zNear = 0.01
        camera.fieldOfView = 60
        
        globeView.cameraNode.camera = camera
        globeView.cameraNode.position = SCNVector3(x: 0, y: 0, z: 2.5)
        
        globeView.sceneView.scene?.rootNode.addChildNode(globeView.cameraNode)
        globeView.sceneView.allowsCameraControl = true
        globeView.sceneView.backgroundColor = UIColor.black
        globeView.sceneView.autoenablesDefaultLighting = true
        globeView.setupTapGesture()
        globeView.updateGlobe(pads: pads, isRotating: isRotating, showLabels: showLabels, selectedCountries: selectedCountries)
        return globeView
    }

    func updateUIView(_ uiView: UIView, context: Context) {
        guard let globeView = uiView as? GlobeSceneView else { return }
        globeView.updateGlobe(pads: pads, isRotating: isRotating, showLabels: showLabels, selectedCountries: selectedCountries)
    }
}

class GlobeSceneView: UIView {
    let sceneView = SCNView()
    let cameraNode = SCNNode()

    private var earthNode: SCNNode = SCNNode()
    private var rotationAnimationKey = "earth-rotation"
    private var hasSetupEarth = false

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

        let scene = SCNScene()
        sceneView.scene = scene
        sceneView.backgroundColor = UIColor.black
        sceneView.autoenablesDefaultLighting = true
        sceneView.allowsCameraControl = true

        cameraNode.camera = SCNCamera()
        cameraNode.position = SCNVector3(0, 0, 2.5)
        scene.rootNode.addChildNode(cameraNode)

        // Removed setupEarth() call here to fix blank screen on initial load
    }

    func setupEarth() {
        // Only create earthNode if it has no geometry yet
        if earthNode.geometry == nil {
            let earthGeometry = SCNSphere(radius: 0.8)
            let material = SCNMaterial()
            material.diffuse.contents = UIImage(named: "earth.jpg")
            earthGeometry.firstMaterial = material

            earthNode.geometry = earthGeometry
        }
        if earthNode.parent == nil {
            sceneView.scene?.rootNode.addChildNode(earthNode)
        }
    }

    func setupTapGesture() {
        let tap = UITapGestureRecognizer(target: self, action: #selector(handleTap(_:)))
        sceneView.addGestureRecognizer(tap)
    }

    @objc private func handleTap(_ gesture: UITapGestureRecognizer) {
        let location = gesture.location(in: sceneView)
        let hits = sceneView.hitTest(location, options: nil)
        if let node = hits.first?.node, let name = node.name {
            let alert = UIAlertController(title: name, message: "Launchpad selected", preferredStyle: .alert)
            alert.addAction(UIAlertAction(title: "OK", style: .default))
            if let vc = sceneView.window?.windowScene?.keyWindow?.rootViewController {
                vc.present(alert, animated: true)
            }
        }
    }

    func updateGlobe(pads: [Pad], isRotating: Bool, showLabels: Bool, selectedCountries: Set<String>) {
        // Ensure earthNode is attached and has geometry
        if !hasSetupEarth {
            setupEarth()
            hasSetupEarth = true
        }

        // 1. Toggle rotation without recreating the globe
        if isRotating {
            if earthNode.animation(forKey: rotationAnimationKey) == nil {
                let rotation = CABasicAnimation(keyPath: "rotation")
                rotation.fromValue = NSValue(scnVector4: SCNVector4(0, 1, 0, 0))
                rotation.toValue = NSValue(scnVector4: SCNVector4(0, 1, 0, Float.pi * 2))
                rotation.duration = 30
                rotation.repeatCount = .infinity
                earthNode.addAnimation(rotation, forKey: rotationAnimationKey)
            }
        } else {
            earthNode.removeAnimation(forKey: rotationAnimationKey, blendOutDuration: 0.2)
        }

        // 2. Clear previous markers (pins and labels) but keep earthNode geometry intact
        earthNode.childNodes.forEach { $0.removeFromParentNode() }

        // 3. Add pins + labels based on visible countries
        for pad in pads {
            if let countryCode = pad.country_code, selectedCountries.contains(countryCode) {
                print("Rendering pad: \(pad.name ?? "Unnamed") at lat: \(pad.latitude), lng: \(pad.longitude)")
                let position = latLngToXYZ(lat: pad.latitude, lng: pad.longitude, radius: 0.8)

                let pin = SCNSphere(radius: 0.01)
                let pinMaterial = SCNMaterial()
                pinMaterial.diffuse.contents = UIColor.systemRed
                pin.firstMaterial = pinMaterial

                let pinNode = SCNNode(geometry: pin)
                pinNode.position = position
                pinNode.name = pad.name ?? "Unknown Pad"

                let pinConstraint = SCNLookAtConstraint(target: cameraNode)
                pinConstraint.isGimbalLockEnabled = true
                pinNode.constraints = [pinConstraint]

                earthNode.addChildNode(pinNode)

                if showLabels {
                    let text = SCNText(string: pad.name, extrusionDepth: 0.05)
                    text.font = UIFont.boldSystemFont(ofSize: 1.5) // Use a larger font
                    text.flatness = 0.1
                    text.firstMaterial?.diffuse.contents = UIColor.white
                    text.firstMaterial?.isDoubleSided = true

                    let textNode = SCNNode(geometry: text)
                    textNode.scale = SCNVector3(0.0025, 0.0025, 0.0025)

                    // Offset the label to float above the pin
                    textNode.position = SCNVector3(position.x + 0.02, position.y + 0.025, position.z + 0.02)

                    // Ensure the label always faces the camera
                    let labelConstraint = SCNBillboardConstraint()
                    labelConstraint.freeAxes = .all
                    textNode.constraints = [labelConstraint]

                    earthNode.addChildNode(textNode)
                }
            }
        }
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

