(call target: (identifier) @keyword
  (arguments
    [(call target: (identifier))
     (identifier)] @func_name)) @identifier
(#match? @keyword "^(def|defp)$")
