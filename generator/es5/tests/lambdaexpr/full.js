$module('full', function () {
  var $instance = this;
  $instance.DoSomething = function () {
    var $this = this;
    var $state = {
      current: 0,
      returnValue: null,
    };
    $state.next = function ($callback) {
      try {
        while (true) {
          switch ($state.current) {
            case 0:
              function (firstParam, secondParam) {
                var $state = {
                  current: 0,
                  returnValue: null,
                };
                $state.next = function ($callback) {
                  try {
                    while (true) {
                      switch ($state.current) {
                        case 0:
                          1234;
                          $state.returnValue = 4567;
                          $state.current = -1;
                          $callback($state);
                          return;
                      }
                    }
                  } catch (e) {
                    $state.error = e;
                    $state.current = -1;
                    $callback($state);
                  }
                };
                return $promise.build($state);
              };
              $state.current = -1;
              return;
          }
        }
      } catch (e) {
        $state.error = e;
        $state.current = -1;
        $callback($state);
      }
    };
    return $promise.build($state);
  };
});
