import * as React from 'react'

export const useForm = (initialValues:any, callback:any, validate:any) => {
	const [values, setValues] = React.useState<any>(initialValues);
	const [errors, setErrors] = React.useState<any>({});
	const [isSubmitting, setIsSubmitting] = React.useState<boolean>(false);

	React.useEffect(
		() => {
			if (Object.keys(errors).length === 0 && isSubmitting) {
				callback();
			}
		},
		[errors],
	)

	const handleSubmit = (event:any) => {
		if (event) event.preventDefault()
		// Only validate if the validate function is used
		if (validate) {
            setErrors(validate(values))
        }
        setIsSubmitting(true);
	}

	const handleChange = (event:any) => {
		event.persist()
		setValues((values:any) => ({
			...values,
			[event.target.name]: event.target.value,
		}))
	}

	return {
		handleChange,
		handleSubmit,
		values,
        errors,
	}
}

