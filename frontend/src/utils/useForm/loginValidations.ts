export const validate = (values: any) => {
    let errors: any = {};
    if (!values.username) {
        errors.username = "Username is required";
    }
    return errors;
};