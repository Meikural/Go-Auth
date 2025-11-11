#!/usr/bin/env python3

import requests
import json
import base64
import os
from pathlib import Path
from datetime import datetime
from dotenv import load_dotenv

# Load .env from root directory
env_path = Path(__file__).parent.parent / ".env"
load_dotenv(env_path)

# Configuration from .env
AUTH_SERVICE_URL = os.getenv("AUTH_SERVICE_URL", "http://localhost:8080")
SUPER_ADMIN_EMAIL = os.getenv("SUPER_ADMIN_EMAIL", "superadmin@web.com")
SUPER_ADMIN_PASSWORD = os.getenv("SUPER_ADMIN_PASSWORD", "superadminpass123")

# Parse ROLES from .env
try:
    roles_str = os.getenv("ROLES", '["Super Admin", "User"]')
    ROLES = json.loads(roles_str)
except json.JSONDecodeError:
    ROLES = ["Super Admin", "User"]

# Color codes for output
GREEN = '\033[92m'
RED = '\033[91m'
BLUE = '\033[94m'
YELLOW = '\033[93m'
RESET = '\033[0m'

# Test results tracking
test_results = {}

def print_header(text):
    print(f"\n{BLUE}{'='*60}")
    print(f"{text}")
    print(f"{'='*60}{RESET}\n")

def print_success(text):
    print(f"{GREEN}✓ {text}{RESET}")

def print_error(text):
    print(f"{RED}✗ {text}{RESET}")

def print_info(text):
    print(f"{YELLOW}ℹ {text}{RESET}")

def decode_jwt(token):
    """Decode JWT token to view claims"""
    try:
        parts = token.split('.')
        if len(parts) != 3:
            return None
        
        payload = parts[1]
        padding = 4 - len(payload) % 4
        if padding != 4:
            payload += '=' * padding
        
        decoded = base64.urlsafe_b64decode(payload)
        return json.loads(decoded)
    except Exception as e:
        print_error(f"Failed to decode token: {e}")
        return None

def test_health_check():
    print_header("Testing Health Check")
    try:
        response = requests.get(f"{AUTH_SERVICE_URL}/health")
        if response.status_code == 200:
            print_success("Health check passed")
            print(f"Response: {response.json()}")
            return True
        else:
            print_error(f"Health check failed with status {response.status_code}")
            return False
    except Exception as e:
        print_error(f"Health check error: {e}")
        return False

def test_super_admin_login():
    print_header("Testing Super Admin Login")
    try:
        payload = {
            "email": SUPER_ADMIN_EMAIL,
            "password": SUPER_ADMIN_PASSWORD
        }
        response = requests.post(f"{AUTH_SERVICE_URL}/login", json=payload)
        
        if response.status_code == 200:
            data = response.json()
            access_token = data.get("access_token")
            refresh_token = data.get("refresh_token")
            user_role = data['user']['role']
            user_id = data['user']['id']
            
            print_success("Super admin login successful")
            print(f"Email: {data['user']['email']}")
            print(f"Role: {user_role}")
            
            claims = decode_jwt(access_token)
            if claims:
                print(f"\nToken Claims:")
                print(f"  User ID: {claims.get('user_id')}")
                print(f"  Username: {claims.get('username')}")
                print(f"  Role: {claims.get('role')}")
                print(f"  Token Type: {claims.get('token_type')}")
            
            return {
                "access_token": access_token,
                "refresh_token": refresh_token,
                "role": user_role,
                "id": user_id
            }
        else:
            print_error(f"Super admin login failed: {response.text}")
            return None
    except Exception as e:
        print_error(f"Super admin login error: {e}")
        return None

def test_user_registration(username, email, password):
    print_header(f"Testing User Registration - {username}")
    try:
        payload = {
            "username": username,
            "email": email,
            "password": password
        }
        response = requests.post(f"{AUTH_SERVICE_URL}/register", json=payload)
        
        if response.status_code == 201:
            data = response.json()
            access_token = data.get("access_token")
            refresh_token = data.get("refresh_token")
            user_role = data['user']['role']
            user_id = data['user']['id']
            
            print_success("User registration successful")
            print(f"Username: {data['user']['username']}")
            print(f"Email: {data['user']['email']}")
            print(f"Role: {user_role}")
            
            claims = decode_jwt(access_token)
            if claims:
                print(f"\nToken Claims:")
                print(f"  User ID: {claims.get('user_id')}")
                print(f"  Role: {claims.get('role')}")
            
            return {
                "access_token": access_token,
                "refresh_token": refresh_token,
                "role": user_role,
                "email": email,
                "password": password,
                "id": user_id
            }
        else:
            print_error(f"User registration failed: {response.text}")
            return None
    except Exception as e:
        print_error(f"User registration error: {e}")
        return None

def test_update_user_role(super_admin_token, user_id, new_role, role_name):
    print_header(f"Testing Update User Role - {role_name}")
    try:
        headers = {
            "Authorization": f"Bearer {super_admin_token}"
        }
        payload = {
            "role": new_role
        }
        response = requests.put(
            f"{AUTH_SERVICE_URL}/admin/users/{user_id}/role",
            json=payload,
            headers=headers
        )
        
        if response.status_code == 200:
            data = response.json()
            updated_user = data.get("user")
            print_success(f"User role updated successfully")
            print(f"User ID: {updated_user['id']}")
            print(f"Username: {updated_user['username']}")
            print(f"New Role: {updated_user['role']}")
            return True
        else:
            print_error(f"Update user role failed: {response.text}")
            return False
    except Exception as e:
        print_error(f"Update user role error: {e}")
        return False

def test_user_login(email, password, expected_role):
    print_header(f"Testing User Login - {email}")
    try:
        payload = {
            "email": email,
            "password": password
        }
        response = requests.post(f"{AUTH_SERVICE_URL}/login", json=payload)
        
        if response.status_code == 200:
            data = response.json()
            access_token = data.get("access_token")
            user_role = data['user']['role']
            
            print_success("User login successful")
            print(f"Email: {data['user']['email']}")
            print(f"Role: {user_role}")
            
            if user_role.lower() != expected_role.lower():
                print_error(f"Expected role {expected_role}, got {user_role}")
            else:
                print_success(f"Role verified: {user_role}")
            
            claims = decode_jwt(access_token)
            if claims:
                print(f"\nToken Claims:")
                print(f"  Role: {claims.get('role')}")
            
            return access_token
        else:
            print_error(f"User login failed: {response.text}")
            return None
    except Exception as e:
        print_error(f"User login error: {e}")
        return None

def test_get_profile(access_token, role_name):
    print_header(f"Testing Get Profile - {role_name}")
    try:
        headers = {
            "Authorization": f"Bearer {access_token}"
        }
        response = requests.get(
            f"{AUTH_SERVICE_URL}/profile",
            headers=headers
        )
        
        if response.status_code == 200:
            data = response.json()
            print_success("Profile retrieved successfully")
            print(f"Username: {data['username']}")
            print(f"Email: {data['email']}")
            print(f"Role: {data['role']}")
            return True
        else:
            print_error(f"Get profile failed: {response.text}")
            return False
    except Exception as e:
        print_error(f"Get profile error: {e}")
        return False

def test_change_password(access_token, old_password, new_password, role_name):
    print_header(f"Testing Change Password - {role_name}")
    try:
        headers = {
            "Authorization": f"Bearer {access_token}"
        }
        payload = {
            "old_password": old_password,
            "new_password": new_password
        }
        response = requests.post(
            f"{AUTH_SERVICE_URL}/change-password",
            json=payload,
            headers=headers
        )
        
        if response.status_code == 200:
            print_success("Password changed successfully")
            print(f"Response: {response.json()}")
            return True
        else:
            print_error(f"Change password failed: {response.text}")
            return False
    except Exception as e:
        print_error(f"Change password error: {e}")
        return False

def test_refresh_token(refresh_token, role_name):
    print_header(f"Testing Token Refresh - {role_name}")
    try:
        payload = {
            "refresh_token": refresh_token
        }
        response = requests.post(f"{AUTH_SERVICE_URL}/refresh", json=payload)
        
        if response.status_code == 200:
            data = response.json()
            new_access_token = data.get("access_token")
            print_success("Token refresh successful")
            
            claims = decode_jwt(new_access_token)
            if claims:
                print(f"\nNew Token Claims:")
                print(f"  User ID: {claims.get('user_id')}")
                print(f"  Role: {claims.get('role')}")
            
            return new_access_token
        else:
            print_error(f"Token refresh failed: {response.text}")
            return None
    except Exception as e:
        print_error(f"Token refresh error: {e}")
        return None

def test_get_all_users(super_admin_token):
    print_header("Testing Get All Users (Admin Endpoint)")
    try:
        headers = {
            "Authorization": f"Bearer {super_admin_token}"
        }
        response = requests.get(
            f"{AUTH_SERVICE_URL}/admin/users",
            headers=headers
        )
        
        if response.status_code == 200:
            data = response.json()
            total = data.get("total")
            users = data.get("users", [])
            print_success(f"Retrieved {total} users")
            
            print(f"\n{'ID':<5} {'Username':<20} {'Email':<30} {'Role':<15}")
            print("-" * 75)
            for user in users:
                print(f"{user['id']:<5} {user['username']:<20} {user['email']:<30} {user['role']:<15}")
            
            return True
        else:
            print_error(f"Get all users failed: {response.text}")
            return False
    except Exception as e:
        print_error(f"Get all users error: {e}")
        return False

def test_unauthorized_access_to_admin_endpoint(user_token, username):
    print_header(f"Testing Unauthorized Access - {username} accessing admin endpoint")
    try:
        headers = {
            "Authorization": f"Bearer {user_token}"
        }
        response = requests.get(
            f"{AUTH_SERVICE_URL}/admin/users",
            headers=headers
        )
        
        if response.status_code == 403:
            print_success(f"Access correctly denied (403 Forbidden)")
            data = response.json()
            print(f"Response: {data}")
            return True
        else:
            print_error(f"Expected 403 Forbidden, got {response.status_code}")
            return False
    except Exception as e:
        print_error(f"Unauthorized access test error: {e}")
        return False

def test_unauthorized_role_update(user_token, target_user_id, new_role, username):
    print_header(f"Testing Unauthorized Role Update - {username} trying to update role")
    try:
        headers = {
            "Authorization": f"Bearer {user_token}"
        }
        payload = {
            "role": new_role
        }
        response = requests.put(
            f"{AUTH_SERVICE_URL}/admin/users/{target_user_id}/role",
            json=payload,
            headers=headers
        )
        
        if response.status_code == 403:
            print_success(f"Role update correctly denied (403 Forbidden)")
            data = response.json()
            print(f"Response: {data}")
            return True
        else:
            print_error(f"Expected 403 Forbidden, got {response.status_code}")
            return False
    except Exception as e:
        print_error(f"Unauthorized role update test error: {e}")
        return False

def print_test_results_summary():
    print_header("Test Results Summary")
    
    if not test_results:
        print_info("No tests were run")
        return
    
    print(f"{'Role':<15} {'Tests':<10} {'Passed':<10} {'Failed':<10}")
    print("-" * 50)
    
    total_passed = 0
    total_failed = 0
    
    for role, results in test_results.items():
        passed = sum(1 for r in results.values() if r)
        failed = len(results) - passed
        total_passed += passed
        total_failed += failed
        
        status = GREEN if failed == 0 else RED
        print(f"{role:<15} {len(results):<10} {status}{passed:<10}{RESET} {failed:<10}")
    
    print("-" * 50)
    print(f"{'TOTAL':<15} {total_passed + total_failed:<10} {GREEN}{total_passed:<10}{RESET} {RED}{total_failed:<10}{RESET}")
    
    if total_failed == 0:
        print_success("All tests passed!")
    else:
        print_error(f"{total_failed} tests failed")

def main():
    print(f"{BLUE}{'='*60}")
    print("Auth Service - Multi-Role Test Suite")
    print(f"{'='*60}{RESET}")
    print(f"Auth Service URL: {AUTH_SERVICE_URL}")
    print(f"Roles to Test: {', '.join(ROLES)}")
    print(f"Timestamp: {datetime.now()}\n")

    # Test 1: Health Check
    if not test_health_check():
        print_error("Auth service is not running!")
        return

    # Test 2: Super Admin Login
    super_admin_data = test_super_admin_login()
    if not super_admin_data:
        print_error("Super admin login failed!")
        return

    test_results["Super Admin"] = {
        "login": True,
        "get_profile": False,
        "change_password": False,
        "token_refresh": False
    }

    # Test 3: Super Admin - Get Profile
    if super_admin_data:
        result = test_get_profile(super_admin_data["access_token"], "Super Admin")
        test_results["Super Admin"]["get_profile"] = result

    # Test 4: Super Admin - Change Password
    if super_admin_data:
        result = test_change_password(
            super_admin_data["access_token"],
            SUPER_ADMIN_PASSWORD,
            "newadminpass123",
            "Super Admin"
        )
        test_results["Super Admin"]["change_password"] = result

    # Test 5: Super Admin - Refresh Token
    if super_admin_data:
        result = test_refresh_token(super_admin_data["refresh_token"], "Super Admin")
        test_results["Super Admin"]["token_refresh"] = result is not None

    # Test each role by creating a user and updating role
    for idx, role in enumerate(ROLES):
        if role.lower() == "super admin":
            continue
        
        test_results[role] = {
            "register": False,
            "update_role": False,
            "login": False,
            "get_profile": False,
            "change_password": False,
            "token_refresh": False
        }
        
        # Register user (will have default User role)
        username = f"test{role.lower().replace(' ', '')}_{idx}"
        email = f"test{role.lower().replace(' ', '')}_{idx}@example.com"
        password = f"testpass{idx}123"
        
        user_data = test_user_registration(username, email, password)
        if user_data:
            test_results[role]["register"] = True
            user_id = user_data.get("id")
            
            # Update user role using super admin endpoint
            if super_admin_data and user_id:
                result = test_update_user_role(
                    super_admin_data["access_token"],
                    user_id,
                    role,
                    role
                )
                test_results[role]["update_role"] = result
            
            # Login again to get new token with updated role
            login_token = test_user_login(email, password, role)
            if login_token:
                test_results[role]["login"] = True
                
                # Test get profile with new role
                result = test_get_profile(login_token, role)
                test_results[role]["get_profile"] = result
                
                # Test change password
                new_password = f"newpass{idx}123"
                result = test_change_password(
                    login_token,
                    password,
                    new_password,
                    role
                )
                test_results[role]["change_password"] = result
                
                # Test token refresh
                result = test_refresh_token(user_data["refresh_token"], role)
                test_results[role]["token_refresh"] = result is not None

    # Test Get All Users endpoint
    if super_admin_data:
        test_get_all_users(super_admin_data["access_token"])

    # Test unauthorized access to admin endpoints
    print_header("Testing Authorization & Access Control")
    
    # Store first non-super-admin user for unauthorized tests
    test_user_token = None
    test_user_id = None
    test_username = None
    
    for idx, role in enumerate(ROLES):
        if role.lower() == "super admin":
            continue
        
        # Get a token from one of the non-admin users
        username = f"test{role.lower().replace(' ', '')}_{idx}"
        email = f"test{role.lower().replace(' ', '')}_{idx}@example.com"
        password = f"newpass{idx}123"
        
        # Try to login to get token
        try:
            payload = {
                "email": email,
                "password": password
            }
            response = requests.post(f"{AUTH_SERVICE_URL}/login", json=payload)
            if response.status_code == 200:
                data = response.json()
                test_user_token = data.get("access_token")
                test_user_id = data['user']['id']
                test_username = username
                break
        except:
            pass
    
    if test_user_token and test_user_id:
        # Test 1: Non-admin user trying to access /admin/users
        test_unauthorized_access_to_admin_endpoint(test_user_token, test_username)
        
        # Test 2: Non-admin user trying to update another user's role
        if super_admin_data:
            test_unauthorized_role_update(
                test_user_token,
                super_admin_data["id"],
                "User",
                test_username
            )

    # Print summary
    print_test_results_summary()

if __name__ == "__main__":
    main()