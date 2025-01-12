extern crate similar;

use similar::TextDiff;
use std::ffi::{CStr, CString};

#[no_mangle]
pub extern "C" fn hello(name: *const libc::c_char) {
    let name_cstr = unsafe { CStr::from_ptr(name) };
    let name = name_cstr.to_str().unwrap();
    println!("Hello {}!", name);
}

#[no_mangle]
pub extern "C" fn diff(prev: *const libc::c_char, newtt: *const libc::c_char) -> *mut libc::c_char {
    let message_cstr = unsafe { CStr::from_ptr(prev) };
    let message = message_cstr.to_str().unwrap();

    let new_message_cstr = unsafe { CStr::from_ptr(newtt) };
    let new_message = new_message_cstr.to_str().unwrap();

    let diff = TextDiff::from_lines(message, new_message);

    CString::new(diff.unified_diff().to_string())
        .unwrap()
        .into_raw()
}
