import { ApolloClient, InMemoryCache, createHttpLink, gql } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import type { 
  User, 
  AuthPayload, 
  LoginRequest, 
  CreateUserRequest 
} from '../types';

const httpLink = createHttpLink({
  uri: 'http://localhost:8080/user/query',
});

const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    }
  };
});

export const apolloClient = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
});

// GraphQL Queries and Mutations
export const FETCH_USERS = gql`
  query FetchUsers {
    fetchUsers {
      id
      username
      email
      role
    }
  }
`;

export const CREATE_USER = gql`
  mutation CreateUser($username: String!, $email: String!, $password: String!, $role: String!) {
    createUser(username: $username, email: $email, password: $password, role: $role) {
      id
      username
      email
      role
    }
  }
`;

export const LOGIN_USER = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      token
      user {
        id
        username
        email
        role
      }
    }
  }
`;

export const LOGOUT_USER = gql`
  mutation Logout {
    logout
  }
`;

// User Service Functions
export const userService = {
  // GraphQL functions
  async fetchUsers(): Promise<User[]> {
    const result = await apolloClient.query({
      query: FETCH_USERS,
    });
    return result.data.fetchUsers;
  },

  async createUser(userData: CreateUserRequest): Promise<User> {
    const result = await apolloClient.mutate({
      mutation: CREATE_USER,
      variables: userData,
    });
    return result.data.createUser;
  },

  async login(credentials: LoginRequest): Promise<AuthPayload> {
    const result = await apolloClient.mutate({
      mutation: LOGIN_USER,
      variables: credentials,
    });
    const authPayload = result.data.login;
    
    // Store token in localStorage
    localStorage.setItem('token', authPayload.token);
    localStorage.setItem('user', JSON.stringify(authPayload.user));
    
    return authPayload;
  },

  async logout(): Promise<boolean> {
    const result = await apolloClient.mutate({
      mutation: LOGOUT_USER,
    });
    
    // Clear local storage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    return result.data.logout;
  },
};
