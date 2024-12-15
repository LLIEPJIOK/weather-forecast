import axios from "axios"
import Cookies from "js-cookie"
import React, { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import "./UpdateWeather.css"

const UpdateWeather = () => {
	const { id } = useParams() // Getting the ID from the URL params
	const navigate = useNavigate()
	const [observation, setObservation] = useState({
		timestamp: "",
		temperature: "",
		humidity: "",
		pressure: "",
		wind_speed: "",
		city: "",
		country: "",
		weather_status: "",
	})

	useEffect(() => {
		var userRole = Cookies.get("X-User-Role")
		if (userRole !== "admin") {
			navigate("/login")
		}

		axios
			.get(`http://localhost:8080/weather/${id}`, {
				withCredentials: true,
			})
			.then(response => {
				const timestamp = response.data.timestamp
					? new Date(response.data.timestamp).toISOString().split("T")[0]
					: ""
				setObservation({ ...response.data, timestamp })
			})
			.catch(error => console.error("Error fetching observation:", error))
	}, [id])

	// Handling changes to input fields
	const handleChange = e => {
		const { name, value } = e.target
		setObservation({ ...observation, [name]: value })
	}

	// Handling form submission to update the observation
	const handleSubmit = e => {
		e.preventDefault()
		const formattedObservation = {
			...observation,
			timestamp: new Date(observation.timestamp).toISOString(), // Приводим к RFC3339
			temperature: parseFloat(observation.temperature), // Преобразуем в число
			humidity: parseFloat(observation.humidity), // Преобразуем в число
			pressure: parseFloat(observation.pressure), // Преобразуем в число
			wind_speed: parseFloat(observation.wind_speed), // Преобразуем в число
			weather_status: observation.weather_status,
		}
		axios
			.put(`http://localhost:8080/weather/${id}`, formattedObservation)
			.then(() => navigate(`/details/${id}`)) // Redirect to the main page after successful update
			.catch(error => console.error("Error updating observation:", error))
	}

	return (
		<div className="update-container">
			<h1 className="update-title">Edit Weather Observation</h1>
			<form onSubmit={handleSubmit} className="update-form">
				<div className="form-group">
					<label>Date:</label>
					<input
						name="timestamp"
						type="date"
						value={observation.timestamp}
						onChange={handleChange}
						className="form-control"
					/>
				</div>
				<div className="form-group">
					<label>Temperature (°C):</label>
					<input
						name="temperature"
						type="number"
						value={observation.temperature}
						onChange={handleChange}
						className="form-control"
					/>
				</div>
				<div className="form-group">
					<label>Humidity (%):</label>
					<input
						name="humidity"
						type="number"
						value={observation.humidity}
						onChange={handleChange}
						className="form-control"
					/>
				</div>
				<div className="form-group">
					<label>Pressure (hPa):</label>
					<input
						name="pressure"
						type="number"
						value={observation.pressure}
						onChange={handleChange}
						className="form-control"
					/>
				</div>
				<div className="form-group">
					<label>Wind Speed (m/s):</label>
					<input
						name="wind_speed"
						type="number"
						value={observation.wind_speed}
						onChange={handleChange}
						className="form-control"
					/>
				</div>
				<div className="form-group">
					<label>City:</label>
					<input
						name="city"
						type="text"
						value={observation.city}
						onChange={handleChange}
						className="form-control"
					/>
				</div>
				<div className="form-group">
					<label>Country:</label>
					<input
						name="country"
						type="text"
						value={observation.country}
						onChange={handleChange}
						className="form-control"
					/>
				</div>
				<div className="form-group">
					<label>Weather Status:</label>
					<input
						name="weather_status"
						type="text"
						value={observation.weather_status}
						onChange={handleChange}
						className="form-control"
					/>
				</div>

				<div className="form-actions">
					<button type="submit" className="save-button">
						Save
					</button>
					<button
						type="button"
						className="cancel-button"
						onClick={() => navigate(`/details/${id}`)}>
						Cancel
					</button>
				</div>
			</form>
		</div>
	)
}

export default UpdateWeather
