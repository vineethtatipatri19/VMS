
import React, {useState} from 'react';
import { View, Text, TextInput, Button, Alert } from 'react-native';
import axios from 'axios';

export default function LoginScreen({navigation}) {
  const [email,setEmail]=useState('');
  const [password,setPassword]=useState('');
  const login = async () => {
    try {
      const res = await axios.post('http://YOUR_BACKEND_URL/api/v1/login',{email,password});
      const token = res.data.token;
      // store token using AsyncStorage or context (left as TODO)
      Alert.alert('Logged in','Token received (store securely)');
      navigation.replace('Dashboard');
    } catch(e) {
      Alert.alert('Login failed', e.response?.data || e.message);
    }
  };
  return (
    <View style={{flex:1,padding:16}}>
      <Text>Login</Text>
      <TextInput placeholder="Email" value={email} onChangeText={setEmail} style={{borderWidth:1,marginBottom:8,padding:8}} />
      <TextInput placeholder="Password" value={password} onChangeText={setPassword} secureTextEntry style={{borderWidth:1,marginBottom:8,padding:8}} />
      <Button title="Login" onPress={login} />
      <Button title="Register" onPress={()=>navigation.navigate('Register')} />
    </View>
  );
}
