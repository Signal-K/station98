//
//  AuthView.swift
//  Station98
//
//  Created by Liam Arbuckle on 7/9/2025.
//

import SwiftUI

enum FocusedField {
    case email
    case password
}

struct LoginView: View {
    @State private var emailText = ""
    @State private var isValidEmail = true
    
    @State private var passwordText = ""
    @State private var isValidPassword = true
    
    var canProceed: Bool {
        Validator.validateEmail(emailText) && Validator.validatePassword(passwordText)
    }
    
    @FocusState private var focusedField: FocusedField?
    
    var body: some View {
        NavigationStack {
            VStack {
                EmaiLTextField(emailText: $emailText, isValidEmail: $isValidEmail)
                
                PasswordTextField(passwordText: $passwordText, isValidPassword: $isValidPassword, validatePassword: Validator.validatePassword, errorText: "Your password is not valid", placeholder: "Password")
                
                HStack {
                    Spacer()
                    Button {
                        
                    } label: {
                        Text("Forgot your password?")
                            .foregroundColor(.red)
                            .font(.system(size: 14, weight: .semibold))
                    }
                    .padding(.trailing)
                    .frame(maxWidth: .infinity)
                }
                
                Button {
                    
                } label: {
                    Text("Sign in")
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
            }
        }
    }
}

struct LoginView_Previews: PreviewProvider {
    static var previews: some View {
        LoginView()
    }
}

struct BottomView: View {
    var googleAction: () -> Void
    var facebookAction: () -> Void
    var appleAction: () -> Void
    
    
    var body: some View {
        VStack {
            Text("Or continue with")
                .font(.system(size: 14, weight: .semibold))
                .foregroundColor(.purple)
                .padding(.bottom)
            
            HStack {
                Button {
                    
                } label: {
                    Text("G+")
                }
                .iconButtonStyle
                Button {
                    
                } label: {
                    Text("FB")
                }
                .iconButtonStyle
                Button {
                    
                } label: {
                    Text("🍎")
                }
                .iconButtonStyle
            }
        }
    }
}

extension View {
    var iconButtonStyle: some View {
        self
            .padding()
            .background(.gray)
            .cornerRadius(8)
    }
}

struct EmaiLTextField: View {
    @Binding var emailText: String
    @Binding var isValidEmail: Bool
    
    @FocusState var focusedField: FocusedField?
    
    var body: some View {
        VStack {
            Text("Login")
                .font(.system(size: 30, weight: .bold))
                .foregroundColor(.blue)
                .padding(.bottom)
            
            TextField("Email", text: $emailText)
                .focused($focusedField, equals: .email)
                .padding()
                .background(.gray)
                .cornerRadius(12)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .stroke(.gray, lineWidth: 3)
                )
                .padding(.horizontal)
                .onChange(of: emailText) { newValue in
                    isValidEmail = Validator.validateEmail(newValue)
                }
            if !isValidEmail {
                HStack {
                    Text("Your email is not valid")
                        .foregroundColor(.red)
                        .padding(.leading)
                    Spacer()
                }
            }
        }
    }
}

struct PasswordTextField: View {
    @Binding var passwordText: String
    @Binding var isValidPassword: Bool
    
    let validatePassword: (String) -> Bool
    
    let errorText: String
    let placeholder: String
    
    @FocusState var focusedField: FocusedField?
    
    var body: some View {
        VStack {
            TextField(placeholder, text: $passwordText)
                .focused($focusedField, equals: .password)
                .padding()
                .background(.white)
                .cornerRadius(12)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .stroke(.gray, lineWidth: 3)
                )
                .padding(.horizontal)
                .onChange(of: passwordText) { newValue in
                    isValidPassword = validatePassword(newValue)
                }
            if !isValidPassword {
                HStack {
                    Text(errorText)
                        .foregroundColor(.red)
                        .padding(.leading)
                    Spacer()
                }
            }
        }
    }
}
