
import React from 'react';
import { View, Text, Button } from 'react-native';

export default function HomeScreen({navigation}) {
  return (
    <View style={{flex:1,alignItems:'center',justifyContent:'center'}}>
      <Text>PGVMS Mobile</Text>
      <Button title="Inventory" onPress={()=>navigation.navigate('Inventory')} />
    </View>
  );
}
