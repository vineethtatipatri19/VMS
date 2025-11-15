
import React, {useEffect, useState} from 'react';
import { View, Text, FlatList } from 'react-native';
import axios from 'axios';

export default function InventoryScreen(){
  const [items,setItems]=useState([]);
  useEffect(()=>{
    axios.get('http://YOUR_BACKEND_URL/api/v1/inventory').then(r=>setItems(r.data)).catch(e=>console.log(e))
  },[]);
  return (
    <View style={{flex:1,padding:16}}>
      <Text>Inventory</Text>
      <FlatList data={items} keyExtractor={i=>i.id} renderItem={({item})=>(
        <View style={{padding:8,borderBottomWidth:1}}>
          <Text>{item.name} — {item.quantity} {item.unit} (exp: {item.expiryDate})</Text>
        </View>
      )} />
    </View>
  );
}
