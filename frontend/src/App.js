import React from "react"
import { Route, BrowserRouter as Router, Routes } from "react-router-dom"
import AddWeather from "./components/AddWeather"
import LoginPage from "./components/Login"
import Navbar from "./components/Navbar"
import UpdateWeather from "./components/UpdateWeather"
import WeatherDetails from "./components/WeatherDetails"
import WeatherList from "./components/WeatherList"

const App = () => {
	return (
		<Router>
			<Navbar />
			<div style={{ padding: "20px" }}>
				<Routes>
					<Route path="/" element={<WeatherList />} />
					<Route path="/add" element={<AddWeather />} />
					<Route path="/details/:id" element={<WeatherDetails />} />
					<Route path="/update/:id" element={<UpdateWeather />} />
					<Route path="/login" element={<LoginPage />} />
				</Routes>
			</div>
		</Router>
	)
}

export default App
