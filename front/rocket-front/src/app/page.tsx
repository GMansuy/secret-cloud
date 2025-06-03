"use client";

import { redirect } from "next/navigation";
import { useState } from "react";
import axios from "axios";

export default function Home() {
    redirect("/login");
    return null;
}