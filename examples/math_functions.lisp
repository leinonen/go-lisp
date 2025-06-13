;;; Mathematical Functions Examples
;;; Demonstrates the new built-in mathematical functions

(println! "=== MATHEMATICAL FUNCTIONS DEMO ===")

;; Basic mathematical functions
(println! "\n--- Basic Math Functions ---")
(println! "sqrt(16) =" (sqrt 16))
(println! "sqrt(2) =" (sqrt 2))
(println! "pow(2, 3) =" (pow 2 3))
(println! "pow(5, 0) =" (pow 5 0))
(println! "abs(-7) =" (abs -7))
(println! "abs(3) =" (abs 3))

;; Rounding functions
(println! "\n--- Rounding Functions ---")
(println! "floor(3.7) =" (floor 3.7))
(println! "ceil(3.2) =" (ceil 3.2))
(println! "round(3.4) =" (round 3.4))
(println! "round(3.6) =" (round 3.6))
(println! "floor(-2.3) =" (floor -2.3))
(println! "ceil(-2.7) =" (ceil -2.7))
(println! "trunc(3.7) =" (trunc 3.7))
(println! "trunc(-3.7) =" (trunc -3.7))

;; Trigonometric functions
(println! "\n--- Trigonometric Functions ---")
(println! "sin(0) =" (sin 0))
(println! "cos(0) =" (cos 0))
(println! "tan(0) =" (tan 0))
(println! "sin(pi/2) ≈" (sin (/ (pi) 2)))
(println! "cos(pi) ≈" (cos (pi)))

;; Inverse trigonometric functions
(println! "\n--- Inverse Trigonometric Functions ---")
(println! "asin(0.5) =" (asin 0.5))
(println! "acos(0.5) =" (acos 0.5))
(println! "atan(1) =" (atan 1))
(println! "atan2(1, 1) =" (atan2 1 1))

;; Hyperbolic functions
(println! "\n--- Hyperbolic Functions ---")
(println! "sinh(0) =" (sinh 0))
(println! "cosh(0) =" (cosh 0))
(println! "tanh(0) =" (tanh 0))
(println! "sinh(1) =" (sinh 1))

;; Angle conversion
(println! "\n--- Angle Conversion ---")
(println! "degrees(pi) =" (degrees (pi)))
(println! "radians(180) =" (radians 180))
(println! "degrees(pi/2) =" (degrees (/ (pi) 2)))

;; Logarithmic and exponential functions
(println! "\n--- Logarithmic and Exponential ---")
(println! "log(e) =" (log (e)))
(println! "log(1) =" (log 1))
(println! "log10(100) =" (log10 100))
(println! "log2(8) =" (log2 8))
(println! "exp(0) =" (exp 0))
(println! "exp(1) =" (exp 1))
(println! "log(exp(5)) =" (log (exp 5)))

;; Min and max functions
(println! "\n--- Min and Max ---")
(println! "min(10, 5) =" (min 10 5))
(println! "max(10, 5) =" (max 10 5))
(println! "min(-3, 2) =" (min -3 2))
(println! "max(-3, 2) =" (max -3 2))

;; Utility functions
(println! "\n--- Utility Functions ---")
(println! "sign(5) =" (sign 5))
(println! "sign(-3) =" (sign -3))
(println! "sign(0) =" (sign 0))
(println! "mod(7, 3) =" (mod 7 3))
(println! "mod(-7, 3) =" (mod -7 3))

;; Mathematical constants
(println! "\n--- Mathematical Constants ---")
(println! "pi =" (pi))
(println! "e =" (e))
(println! "2*pi =" (* 2 (pi)))
(println! "e^pi =" (pow (e) (pi)))

;; Random number generation
(println! "\n--- Random Numbers ---")
(println! "random() (0 to 1) =" (random))
(println! "random() (0 to 1) =" (random))
(println! "random(10) (0 to 9) =" (random 10))
(println! "random(5, 15) (5 to 14) =" (random 5 15))
(println! "random(5, 15) (5 to 14) =" (random 5 15))

;; Complex mathematical expressions
(println! "\n--- Complex Expressions ---")
(println! "Distance formula: sqrt((3-1)^2 + (4-2)^2) =" 
          (sqrt (+ (pow (- 3 1) 2) (pow (- 4 2) 2))))

(println! "Quadratic formula discriminant: b^2 - 4ac where a=1, b=5, c=6")
(def a 1)
(def b 5) 
(def c 6)
(def discriminant (- (pow b 2) (* 4 a c)))
(println! "Discriminant =" discriminant)
(println! "Root 1 =" (/ (+ (- b) (sqrt discriminant)) (* 2 a)))
(println! "Root 2 =" (/ (- (- b) (sqrt discriminant)) (* 2 a)))

;; Practical examples
(println! "\n--- Practical Examples ---")

;; Calculate compound interest: A = P(1 + r/n)^(nt)
(def principal 1000)
(def rate 0.05)     ; 5% annual rate
(def times 12)      ; compounded monthly  
(def years 10)
(def amount (* principal (pow (+ 1 (/ rate times)) (* times years))))
(println! "Compound interest: $1000 at 5% for 10 years =" amount)

;; Calculate area of circle
(def radius 5)
(def area (* (pi) (pow radius 2)))
(println! "Area of circle with radius 5 =" area)

;; Convert between degrees and radians
(def angle-degrees 45)
(def angle-radians (radians angle-degrees))
(println! "45 degrees in radians =" angle-radians)
(println! "Back to degrees =" (degrees angle-radians))

;; Generate a series of random points
(println! "\nRandom points in unit circle:")
(def generate-point (fn []
  (list (- (* 2 (random)) 1) (- (* 2 (random)) 1))))

(def point1 (generate-point))
(def point2 (generate-point))
(def point3 (generate-point))
(println! "Point 1:" point1)
(println! "Point 2:" point2) 
(println! "Point 3:" point3)

;; Statistical functions example
(println! "\n--- Statistical Example ---")
(def numbers (list 1 4 7 2 9 3 8 5 6))
(def n (length numbers))
(def sum (reduce (fn [acc x] (+ acc x)) 0 numbers))
(def mean (/ sum n))
(def min-val (reduce (fn [acc x] (min acc x)) (first numbers) (rest numbers)))
(def max-val (reduce (fn [acc x] (max acc x)) (first numbers) (rest numbers)))

(println! "Numbers:" numbers)
(println! "Count:" n)
(println! "Sum:" sum)
(println! "Mean:" mean)
(println! "Min:" min-val)
(println! "Max:" max-val)

;; Trigonometric identities verification
(println! "\n--- Trigonometric Identities ---")
(def angle (/ (pi) 4))  ; 45 degrees
(def sin-val (sin angle))
(def cos-val (cos angle))
(def tan-val (tan angle))

(println! "Angle (45°) =" angle)
(println! "sin(45°) =" sin-val)
(println! "cos(45°) =" cos-val)  
(println! "tan(45°) =" tan-val)
(println! "sin²(45°) + cos²(45°) =" (+ (pow sin-val 2) (pow cos-val 2)))
(println! "tan(45°) = sin(45°)/cos(45°) =" (/ sin-val cos-val))

(println! "\n=== MATHEMATICAL FUNCTIONS DEMO COMPLETE ===")
