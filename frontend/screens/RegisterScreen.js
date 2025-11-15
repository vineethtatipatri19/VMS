
import React, {useState} from 'react';
import { View, Text, TextInput, Button, Alert } from 'react-native';
import axios from 'axios';

export default function RegisterScreen({navigation}) {
  const [name,setName]=useState('');
  const [email,setEmail]=useState('');
  const [password,setPassword]=useState('');
  const register = async () => {
    try {
      await axios.post('http://YOUR_BACKEND_URL/api/v1/register',{name,email,password});
      Alert.alert('Registered','Please login');
      navigation.goBack();
    } catch(e) {
      Alert.alert('Register failed', e.response?.data || e.message);
    }
  };
  return (
    <View style={{flex:1,padding:16}}>
      <Text>Register</Text>
      <TextInput placeholder="Name" value={name} onChangeText={setName} style={{borderWidth:1,marginBottom:8,padding:8}} />
      <TextInput placeholder="Email" value={email} onChangeText={setEmail} style={{borderWidth:1,marginBottom:8,padding:8}} />
      <TextInput placeholder="Password" value={password} onChangeText={setPassword} secureTextEntry style={{borderWidth:1,marginBottom:8,padding:8}} />
      <Button title="Register" onPress={register} />
    </View>
  );
}
