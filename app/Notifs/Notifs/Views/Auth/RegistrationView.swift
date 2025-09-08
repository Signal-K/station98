//
//  RegistrationView.swift
//  Station98
//
//  Created by Liam Arbuckle on 9/9/2025.
//

import SwiftUI

struct RegistrationView: View {
    @State private var emailText = ""
    @State private var isValidEmail = true
    
    @State private var passwordText = ""
    @State private var confirmPasswordText = ""
    @State private var isValidPassword = true
    @State private var isValidConfirmPassword = true
    
    var canProceed: Bool {
        Validator.validateEmail(emailText) && Validator.validatePassword(passwordText)
    }
    
    @FocusState private var focusedField: FocusedField?
    
    var body: some View {
        NavigationStack {
            VStack {
                Text("Create Account")
                    .font(.system(size: 30, weight: .bold))
                    .foregroundColor(.blue)
                    .padding(.top, 48)
                
                EmaiLTextField(emailText: $emailText, isValidEmail: $isValidEmail)
                
                PasswordTextField(passwordText: $passwordText, isValidPassword: $isValidPassword, validatePassword: Validator.validatePassword, errorText: "Your password is not valid", placeholder: "Password")
                
                PasswordTextField(passwordText: $confirmPasswordText, isValidPassword: $isValidConfirmPassword, validatePassword: { $0 == passwordText }, errorText: "Your passwords do not match", placeholder: "Confirm Password")
                    .padding(.top)
                
                Spacer()
                
                Button {
                    
                } label: {
                    Text("Sign up")
                        .foregroundColor(.white)
                        .font(.system(size: 20, weight: .semibold))
                }
                .padding(.vertical)
                .frame(maxWidth: .infinity)
                .background(.blue)
                .cornerRadius(12)
                .padding(.horizontal)
                .opacity(canProceed ? 1.0 : 0.5)
                .disabled(!canProceed)
                
                
                Button {
                    
                } label: {
                    Text("Create new account")
                        .font(.system(size: 20, weight: .semibold))
                        .foregroundColor(.gray)
                }
                .padding(.vertical)
                .frame(maxWidth: .infinity)
                .cornerRadius(12)
                .padding(.horizontal)
                
                BottomView(googleAction: {}, facebookAction: {}, appleAction: {})
                
                Spacer()
            }
        }
    }
}
