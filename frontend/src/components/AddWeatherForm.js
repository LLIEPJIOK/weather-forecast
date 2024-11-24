import React, { useState } from "react"
import "./AddWeatherForm.css"

const AddWeatherForm = ({ observation, handleChange, handleSubmit }) => {
	// State for errors and success messages
	const [errors, setErrors] = useState({})
	const [message, setMessage] = useState("") // For success or error message

	// Validation function
	const validate = () => {
		let errors = {}
		// Checking required fields
		if (!observation.timestamp) {
			errors.timestamp = "Date is required"
		}
		if (observation.humidity < 0 || observation.humidity > 100) {
			errors.humidity = "Enter a valid humidity"
		}
		if (observation.pressure < 0) {
			errors.pressure = "Enter a valid pressure"
		}
		if (observation.windSpeed < 0) {
			errors.windSpeed = "Enter a valid wind speed"
		}
		if (observation.weatherStatus === "") {
			errors.weatherStatus = "Weather status is required"
		}
		if (observation.city === "") {
			errors.city = "City is required"
		}
		if (observation.country === "") {
			errors.country = "Country is required"
		}

		return errors
	}

	// Form submit handler
	const handleSubmitForm = e => {
		e.preventDefault()

		// Perform validation
		const validationErrors = validate()

		if (Object.keys(validationErrors).length === 0) {
			// If no errors, submit the form
			handleSubmit(e)

			setErrors({})
			setMessage("Observation added successfully!") // Success message
			showMessage(true) // Show success message
		} else {
			// If errors, store them in state
			setErrors(validationErrors)
			setMessage("There was an error. Please check the form.") // Error message
			showMessage(false) // Show error message
		}
	}

	// Function to add error class to input fields
	const getInputClass = field => {
		return errors[field] ? "form-control error" : "form-control"
	}

	// Function to control message visibility
	const showMessage = isSuccess => {
		const messageElement = document.querySelector(".message")
		messageElement.classList.add("show")
		if (!isSuccess) {
			messageElement.classList.add("error")
		}
		// Hide the message after 5 seconds
		setTimeout(() => {
			messageElement.classList.add("hide") // Start fading out
		}, 5000) // Message fades out after 5 seconds

		// After the fade-out completes, remove the 'show' class to completely hide the message
		setTimeout(() => {
			messageElement.classList.remove("show", "hide")
		}, 6000) // 1 second after fade-out completes
	}

	return (
		<div className="form-container">
			<form onSubmit={handleSubmitForm} className="weather-form">
				<h1 className="form-title">Add Weather Observation</h1>

				<div className="form-group">
					<label>Date:</label>
					<input
						name="timestamp"
						type="date" // Changed type to "date" for the calendar
						value={observation.timestamp}
						onChange={handleChange}
						className={getInputClass("timestamp")}
					/>
					{errors.timestamp && (
						<span className="error-message">{errors.timestamp}</span>
					)}
				</div>

				<div className="form-group">
					<label>Temperature (°C):</label>
					<input
						name="temperature"
						type="number"
						value={observation.temperature}
						onChange={handleChange}
						className={getInputClass("temperature")}
					/>
					{errors.temperature && (
						<span className="error-message">{errors.temperature}</span>
					)}
				</div>

				<div className="form-group">
					<label>Humidity (%):</label>
					<input
						name="humidity"
						type="number"
						value={observation.humidity}
						onChange={handleChange}
						className={getInputClass("humidity")}
					/>
					{errors.humidity && (
						<span className="error-message">{errors.humidity}</span>
					)}
				</div>

				<div className="form-group">
					<label>Pressure (hPa):</label>
					<input
						name="pressure"
						type="number"
						value={observation.pressure}
						onChange={handleChange}
						className={getInputClass("pressure")}
					/>
					{errors.pressure && (
						<span className="error-message">{errors.pressure}</span>
					)}
				</div>

				<div className="form-group">
					<label>Wind Speed (m/s):</label>
					<input
						name="windSpeed"
						type="number"
						value={observation.windSpeed}
						onChange={handleChange}
						className={getInputClass("windSpeed")}
					/>
					{errors.windSpeed && (
						<span className="error-message">{errors.windSpeed}</span>
					)}
				</div>

				<div className="form-group">
					<label>City:</label>
					<input
						name="city"
						value={observation.city}
						onChange={handleChange}
						className={getInputClass("city")}
					/>
					{errors.city && <span className="error-message">{errors.city}</span>}
				</div>
				<div className="form-group">
					<label>Country:</label>
					<input
						name="country"
						value={observation.country}
						onChange={handleChange}
						className={getInputClass("country")}
					/>
					{errors.country && (
						<span className="error-message">{errors.country}</span>
					)}
				</div>

				<div className="form-group">
					<label>Weather Status:</label>
					<input
						name="weatherStatus"
						value={observation.weatherStatus}
						onChange={handleChange}
						className={getInputClass("weatherStatus")}
					/>
					{errors.weatherStatus && (
						<span className="error-message">{errors.weatherStatus}</span>
					)}
				</div>

				<button type="submit" className="submit-btn">
					Add
				</button>
			</form>

			{/* Display success/error message */}
			<div className="message">
				{message && <span className="message-text">{message}</span>}
			</div>
		</div>
	)
}

export default AddWeatherForm
