export const validate = (values: any) => {
    let errors: any = {};

    if (!values.name) {
        errors.name = "Name is required";
    } else if (values.name.length < 1) {
        errors.name = "Name must be at least 2 characters";
    }
    if (!values.username) {
        errors.username = "Username is required";
    } else if (values.username.length < 3) {
        errors.username = "Username must be 4 characters long";
    }

    return errors;
};