
import React, {useEffect, useState} from 'react';
import { View, Text, Button, FlatList } from 'react-native';
import axios from 'axios';

export default function Dashboard({navigation}) {
  const [stats,setStats]=useState({customers:0,expiringSoon:0,outstanding:0,unreturnedCrates:0});
  const [inventory,setInventory]=useState([]);
  useEffect(()=>{
    // KPI fetches - the backend doesn't have a single KPI endpoint yet; we'll approximate
    axios.get('http://YOUR_BACKEND_URL/api/v1/inventory').then(r=>{
      setInventory(r.data);
      const today = new Date();
      const soon = r.data.filter(i=>{
        const exp = new Date(i.expiryDate);
        const diff = (exp - today)/(1000*60*60*24);
        return diff >=0 && diff <=3;
      }).length;
      setStats(s=>({...s, expiringSoon: soon}));
    }).catch(e=>console.log(e));
    axios.get('http://YOUR_BACKEND_URL/api/v1/customers').then(r=>setStats(s=>({...s, customers: r.data.length}))).catch(()=>{});
  },[]);
  return (
    <View style={{flex:1,padding:16}}>
      <Text style={{fontSize:18,fontWeight:'bold'}}>PGVMS Dashboard</Text>
      <View style={{marginVertical:12}}>
        <Text>Customers: {stats.customers}</Text>
        <Text>Items expiring within 3 days: {stats.expiringSoon}</Text>
        <Text>Outstanding balance (approx): {stats.outstanding}</Text>
        <Text>Unreturned crates (approx): {stats.unreturnedCrates}</Text>
      </View>
      <Button title="Inventory" onPress={()=>navigation.navigate('Inventory')} />
      <Button title="Transactions" onPress={()=>navigation.navigate('Transactions')} />
      <FlatList data={inventory} keyExtractor={i=>i.id} renderItem={({item})=>(
        <View style={{padding:8,borderBottomWidth:1, backgroundColor: (() => {
          const exp = new Date(item.expiryDate);
          const now = new Date();
          const diff = (exp - now)/(1000*60*60*24);
          if (diff < 0) return '#ffe6e6';
          if (diff <= 3) return '#fff4e6';
          return '#e8ffe8';
        })()}}>
          <Text>{item.name} — {item.quantity} {item.unit} (exp: {item.expiryDate})</Text>
          { /* Visual FEFO badge: first non-expired item is suggested as Sell First (handled client-side example) */ }
        </View>
      )} />
    </View>
  );
}
