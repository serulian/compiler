$module('mini', function () {
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
              function (someParam) {
                var $state = {
                  current: 0,
                  returnValue: null,
                };
                $state.next = function ($callback) {
                  try {
                    while (true) {
                      switch ($state.current) {
                        case 0:
                          $state.returnValue = someParam;
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
