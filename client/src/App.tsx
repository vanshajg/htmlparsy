import { useEffect, useState } from 'react'
import './App.css';
import axios from 'axios';

function App() {
  const [api_data, setApiData] = useState("")
  useEffect(() => { 
    axios.get("/api/ping").then((response) => {
      setApiData(response.data.message)
    })
  }, [])
  return (
    <div className="App">
     <h1>{api_data}</h1>
    </div>
  );
}

export default App;
