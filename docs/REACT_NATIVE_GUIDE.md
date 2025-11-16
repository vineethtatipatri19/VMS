# React Native Mobile App - Implementation Guide

## Overview
This document outlines how to build iOS and Android mobile apps using React Native, leveraging the design system and component architecture from the web app.

---

## Architecture Strategy

### Shared Business Logic
- **API Services**: Reuse `frontend/src/services/api.js` with minimal modifications
- **Authentication Context**: Port `AuthContext.js` logic to React Native AsyncStorage
- **State Management**: Same patterns (Context API or add Redux/Zustand)
- **Utils & Helpers**: 100% reusable JavaScript logic

### UI Components Translation

The web UI components were designed with React Native in mind. Here's the mapping:

| Web Component | React Native Equivalent | Notes |
|--------------|------------------------|-------|
| `<div>` | `<View>` | Direct replacement |
| `<span>`, `<p>` | `<Text>` | All text must be in Text |
| `<button>` | `<TouchableOpacity>` or `<Pressable>` | Use Pressable for modern apps |
| `<input>` | `<TextInput>` | Similar props |
| `<img>` | `<Image>` | Use react-native-fast-image for performance |
| CSS classes | StyleSheet | Convert to StyleSheet.create() |

---

## Project Setup

### 1. Initialize React Native Project
```bash
# Using React Native CLI (recommended for this project)
npx react-native init PGVMSMobile --template react-native-template-typescript

# Or using Expo (easier for rapid development)
npx create-expo-app PGVMSMobile --template
```

### 2. Install Core Dependencies
```bash
# Navigation
npm install @react-navigation/native @react-navigation/native-stack @react-navigation/bottom-tabs
npm install react-native-screens react-native-safe-area-context

# UI & Icons
npm install lucide-react-native
npm install react-native-svg  # Required for lucide icons

# Charts (same library as web!)
npm install react-native-chart-kit
npm install react-native-svg  # Required for charts

# Forms & Input
npm install react-hook-form

# Storage
npm install @react-native-async-storage/async-storage

# HTTP Client
npm install axios

# Date handling
npm install date-fns

# Bottom Sheet (for mobile-friendly modals)
npm install @gorhom/bottom-sheet
npm install react-native-reanimated react-native-gesture-handler
```

---

## Design System Translation

### Color Variables (Shared)
Create `src/styles/colors.ts`:

```typescript
export const colors = {
  primary: '#2563eb',
  primaryDark: '#1e40af',
  primaryLight: '#60a5fa',
  primaryBg: '#eff6ff',
  
  success: '#10b981',
  successDark: '#059669',
  successBg: '#d1fae5',
  
  warning: '#f59e0b',
  warningDark: '#d97706',
  warningBg: '#fef3c7',
  
  danger: '#ef4444',
  dangerDark: '#dc2626',
  dangerBg: '#fee2e2',
  
  gray50: '#f9fafb',
  gray100: '#f3f4f6',
  gray200: '#e5e7eb',
  gray300: '#d1d5db',
  gray500: '#6b7280',
  gray700: '#374151',
  gray900: '#111827',
  
  white: '#ffffff',
  black: '#000000',
};

export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  '2xl': 48,
  '3xl': 64,
};

export const fontSize = {
  xs: 12,
  sm: 14,
  base: 16,
  lg: 18,
  xl: 20,
  '2xl': 24,
  '3xl': 30,
  '4xl': 36,
};

export const radius = {
  sm: 4,
  md: 8,
  lg: 12,
  xl: 16,
  full: 9999,
};
```

### Component Translation Examples

#### Button Component

**Web Version** (`Button.js`):
```jsx
const Button = ({ children, variant, onClick, ... }) => (
  <button className={`btn btn-${variant}`} onClick={onClick}>
    {children}
  </button>
);
```

**React Native Version** (`Button.tsx`):
```typescript
import React from 'react';
import { TouchableOpacity, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { colors, spacing, fontSize } from '../styles/colors';

interface ButtonProps {
  children: React.ReactNode;
  variant?: 'primary' | 'secondary' | 'success' | 'danger';
  onPress?: () => void;
  loading?: boolean;
  disabled?: boolean;
  fullWidth?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  onPress,
  loading = false,
  disabled = false,
  fullWidth = false,
}) => {
  return (
    <TouchableOpacity
      style={[
        styles.button,
        styles[variant],
        fullWidth && styles.fullWidth,
        (disabled || loading) && styles.disabled,
      ]}
      onPress={onPress}
      disabled={disabled || loading}
      activeOpacity={0.7}
    >
      {loading ? (
        <ActivityIndicator color={colors.white} />
      ) : (
        <Text style={styles.text}>{children}</Text>
      )}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    height: 44, // iOS touch target
    paddingHorizontal: spacing.lg,
    borderRadius: radius.md,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  primary: {
    backgroundColor: colors.primary,
  },
  secondary: {
    backgroundColor: colors.gray200,
  },
  success: {
    backgroundColor: colors.success,
  },
  danger: {
    backgroundColor: colors.danger,
  },
  fullWidth: {
    width: '100%',
  },
  disabled: {
    opacity: 0.5,
  },
  text: {
    color: colors.white,
    fontSize: fontSize.base,
    fontWeight: '600',
  },
});
```

#### Card Component

**React Native Version** (`Card.tsx`):
```typescript
import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { colors, spacing, radius } from '../styles/colors';

interface CardProps {
  children: React.ReactNode;
  title?: string;
  subtitle?: string;
}

export const Card: React.FC<CardProps> = ({ children, title, subtitle }) => {
  return (
    <View style={styles.card}>
      {(title || subtitle) && (
        <View style={styles.header}>
          {title && <Text style={styles.title}>{title}</Text>}
          {subtitle && <Text style={styles.subtitle}>{subtitle}</Text>}
        </View>
      )}
      <View style={styles.body}>{children}</View>
    </View>
  );
};

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.white,
    borderRadius: radius.lg,
    padding: spacing.lg,
    marginBottom: spacing.lg,
    shadowColor: colors.black,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3, // Android shadow
  },
  header: {
    marginBottom: spacing.md,
  },
  title: {
    fontSize: fontSize.lg,
    fontWeight: '600',
    color: colors.gray900,
  },
  subtitle: {
    fontSize: fontSize.sm,
    color: colors.gray500,
    marginTop: spacing.xs,
  },
  body: {
    // Content styling
  },
});
```

#### Input Component

**React Native Version** (`Input.tsx`):
```typescript
import React from 'react';
import { View, Text, TextInput, StyleSheet } from 'react-native';
import { colors, spacing, fontSize, radius } from '../styles/colors';

interface InputProps {
  label?: string;
  value: string;
  onChangeText: (text: string) => void;
  placeholder?: string;
  error?: string;
  secureTextEntry?: boolean;
  keyboardType?: 'default' | 'numeric' | 'email-address' | 'phone-pad';
}

export const Input: React.FC<InputProps> = ({
  label,
  value,
  onChangeText,
  placeholder,
  error,
  secureTextEntry,
  keyboardType = 'default',
}) => {
  return (
    <View style={styles.container}>
      {label && <Text style={styles.label}>{label}</Text>}
      <TextInput
        style={[styles.input, error && styles.inputError]}
        value={value}
        onChangeText={onChangeText}
        placeholder={placeholder}
        placeholderTextColor={colors.gray500}
        secureTextEntry={secureTextEntry}
        keyboardType={keyboardType}
        autoCapitalize="none"
      />
      {error && <Text style={styles.errorText}>{error}</Text>}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: spacing.md,
  },
  label: {
    fontSize: fontSize.sm,
    fontWeight: '500',
    color: colors.gray900,
    marginBottom: spacing.xs,
  },
  input: {
    height: 44,
    backgroundColor: colors.white,
    borderWidth: 1,
    borderColor: colors.gray300,
    borderRadius: radius.md,
    paddingHorizontal: spacing.md,
    fontSize: fontSize.base,
    color: colors.gray900,
  },
  inputError: {
    borderColor: colors.danger,
  },
  errorText: {
    fontSize: fontSize.sm,
    color: colors.danger,
    marginTop: spacing.xs,
  },
});
```

---

## Screen Translation Guide

### Dashboard Screen

**Key Changes for Mobile**:
1. **ScrollView**: Wrap entire dashboard in ScrollView
2. **Responsive Grid**: Use Dimensions API or FlexBox
3. **Charts**: Use `react-native-chart-kit` (similar API to Chart.js)
4. **Pull to Refresh**: Add RefreshControl

```typescript
import React, { useState, useEffect } from 'react';
import { ScrollView, RefreshControl, View, Text, StyleSheet, Dimensions } from 'react-native';
import { LineChart } from 'react-native-chart-kit';
import { Card } from '../components/Card';
import { colors, spacing } from '../styles/colors';

const screenWidth = Dimensions.get('window').width;

export const DashboardScreen = () => {
  const [refreshing, setRefreshing] = useState(false);
  const [stats, setStats] = useState(null);

  const onRefresh = async () => {
    setRefreshing(true);
    await fetchData();
    setRefreshing(false);
  };

  const fetchData = async () => {
    // API call
  };

  const chartConfig = {
    backgroundColor: colors.primary,
    backgroundGradientFrom: colors.primary,
    backgroundGradientTo: colors.primaryLight,
    color: (opacity = 1) => `rgba(255, 255, 255, ${opacity})`,
    strokeWidth: 2,
  };

  return (
    <ScrollView
      style={styles.container}
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
      }
    >
      {/* KPI Cards */}
      <View style={styles.kpiGrid}>
        <Card title="Total Customers">
          <Text style={styles.kpiValue}>{stats?.totalCustomers || 0}</Text>
        </Card>
        {/* More KPI cards */}
      </View>

      {/* Chart */}
      <Card title="Sales Trend">
        <LineChart
          data={{
            labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri'],
            datasets: [{ data: [20, 45, 28, 80, 99] }],
          }}
          width={screenWidth - 60}
          height={220}
          chartConfig={chartConfig}
          bezier
        />
      </Card>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.gray50,
    padding: spacing.md,
  },
  kpiGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    justifyContent: 'space-between',
  },
  kpiValue: {
    fontSize: 32,
    fontWeight: 'bold',
    color: colors.gray900,
  },
});
```

---

## Navigation Structure

### React Navigation Setup

```typescript
// App.tsx
import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Home, Package, Users, DollarSign } from 'lucide-react-native';

// Screens
import { LoginScreen } from './screens/LoginScreen';
import { DashboardScreen } from './screens/DashboardScreen';
import { InventoryScreen } from './screens/InventoryScreen';
import { CustomersScreen } from './screens/CustomersScreen';
import { TransactionsScreen } from './screens/TransactionsScreen';

const Stack = createNativeStackNavigator();
const Tab = createBottomTabNavigator();

const TabNavigator = () => (
  <Tab.Navigator
    screenOptions={{
      tabBarActiveTintColor: colors.primary,
      tabBarInactiveTintColor: colors.gray500,
    }}
  >
    <Tab.Screen
      name="Dashboard"
      component={DashboardScreen}
      options={{
        tabBarIcon: ({ color, size }) => <Home color={color} size={size} />,
      }}
    />
    <Tab.Screen
      name="Inventory"
      component={InventoryScreen}
      options={{
        tabBarIcon: ({ color, size }) => <Package color={color} size={size} />,
      }}
    />
    <Tab.Screen
      name="Customers"
      component={CustomersScreen}
      options={{
        tabBarIcon: ({ color, size }) => <Users color={color} size={size} />,
      }}
    />
    <Tab.Screen
      name="Transactions"
      component={TransactionsScreen}
      options={{
        tabBarIcon: ({ color, size }) => <DollarSign color={color} size={size} />,
      }}
    />
  </Tab.Navigator>
);

export default function App() {
  return (
    <NavigationContainer>
      <Stack.Navigator>
        <Stack.Screen name="Login" component={LoginScreen} options={{ headerShown: false }} />
        <Stack.Screen name="Main" component={TabNavigator} options={{ headerShown: false }} />
      </Stack.Navigator>
    </NavigationContainer>
  );
}
```

---

## API Service Reusability

The API service is 95% reusable! Just change axios defaults:

```typescript
// services/api.ts
import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

const API_URL = 'https://bug-free-waddle-7xx5jqgqr72rxvx-8080.app.github.dev/api/v1';

const api = axios.create({
  baseURL: API_URL,
  timeout: 10000,
});

// Request interceptor - add token
api.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - handle errors
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      await AsyncStorage.removeItem('token');
      // Navigate to login
    }
    return Promise.reject(error);
  }
);

// Export same API functions as web!
export const authAPI = {
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData),
};

export const inventoryAPI = {
  getAll: () => api.get('/inventory'),
  create: (data) => api.post('/inventory', data),
  update: (id, data) => api.put(`/inventory/${id}`, data),
  delete: (id) => api.delete(`/inventory/${id}`),
};

// ... same for all other APIs
```

---

## Mobile-Specific Features

### 1. Camera Integration (for scanning barcodes)
```bash
npm install react-native-camera-kit
# or
npm install expo-camera  # if using Expo
```

### 2. Offline Mode
```bash
npm install @react-native-async-storage/async-storage
npm install @tanstack/react-query  # For offline caching
```

### 3. Push Notifications
```bash
npm install @react-native-firebase/messaging  # FCM
# or
npm install expo-notifications  # if using Expo
```

### 4. Biometric Authentication
```bash
npm install react-native-biometrics
```

---

## Platform-Specific Considerations

### iOS
- **Min Version**: iOS 13+
- **Permissions**: Camera, Notifications in Info.plist
- **TouchID/FaceID**: Native support via react-native-biometrics

### Android
- **Min SDK**: 21 (Android 5.0)
- **Permissions**: Camera, Notifications in AndroidManifest.xml
- **Fingerprint**: Native support

---

## Development Workflow

### 1. Run on Simulator
```bash
# iOS
npx react-native run-ios

# Android
npx react-native run-android
```

### 2. Hot Reload
React Native has Fast Refresh enabled by default (like web hot reload)

### 3. Debug
```bash
# Open dev menu
# iOS: Cmd+D
# Android: Cmd+M (Mac) or Ctrl+M (Windows)
```

### 4. Build for Production

**iOS** (requires Mac + Xcode):
```bash
cd ios
pod install
cd ..
npx react-native run-ios --configuration Release
```

**Android**:
```bash
cd android
./gradlew assembleRelease
# APK will be in android/app/build/outputs/apk/release/
```

---

## Code Sharing Strategy

### Recommended Monorepo Structure
```
pgvms/
├── packages/
│   ├── shared/              # Shared business logic
│   │   ├── api/            # API services (100% shared)
│   │   ├── utils/          # Helper functions (100% shared)
│   │   ├── constants/      # Constants (100% shared)
│   │   └── types/          # TypeScript types (100% shared)
│   ├── web/                # React web app (current)
│   └── mobile/             # React Native app (new)
├── backend/                # Go backend (unchanged)
└── infra/                  # Infrastructure (unchanged)
```

### Using Yarn Workspaces or npm Workspaces
```json
// package.json (root)
{
  "name": "pgvms-monorepo",
  "private": true,
  "workspaces": [
    "packages/*"
  ]
}
```

---

## Timeline Estimate

| Phase | Tasks | Duration |
|-------|-------|----------|
| **Phase 1: Setup** | Project init, navigation, UI components | 1 week |
| **Phase 2: Core Screens** | Dashboard, Inventory, Customers, Transactions | 2 weeks |
| **Phase 3: Features** | Camera scanning, offline mode, charts | 1 week |
| **Phase 4: Polish** | Animations, error handling, testing | 1 week |
| **Phase 5: Release** | App store submissions, beta testing | 1 week |
| **Total** | | **6 weeks** |

---

## Resources

### Documentation
- [React Native Docs](https://reactnative.dev/docs/getting-started)
- [React Navigation](https://reactnavigation.org/docs/getting-started)
- [React Native Chart Kit](https://github.com/indiespirit/react-native-chart-kit)

### Tools
- [React Native Debugger](https://github.com/jhen0409/react-native-debugger)
- [Flipper](https://fbflipper.com/) - Advanced debugging
- [Reactotron](https://github.com/infinitered/reactotron) - Inspector

---

## Next Steps

1. **Review this guide** with the development team
2. **Decide on tooling**: React Native CLI vs Expo
3. **Set up development environment**: Xcode (iOS) and Android Studio
4. **Create proof of concept**: Build Login + Dashboard screens
5. **Iterate based on feedback**

---

**Note**: The web app was intentionally designed with mobile compatibility in mind. Component props, state management, and business logic can be reused with minimal changes. Focus on translating the UI layer to React Native components.
