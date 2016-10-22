$module('children', function () {
  var $static = this;
  $static.SimpleFunction = function (props, children) {
    var $result;
    var $temp0;
    var $temp1;
    var counter;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            counter = $t.fastbox(0, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          case 1:
            $temp1 = children;
            $current = 2;
            continue;

          case 2:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            value = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue;
            } else {
              $current = 5;
              continue;
            }
            break;

          case 4:
            counter = $t.fastbox(counter.$wrapped + value.$wrapped, $g.____testlib.basictypes.Boolean);
            $current = 2;
            continue;

          case 5:
            $resolve($t.fastbox(counter.$wrapped == 6, $g.____testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty().then(function ($result1) {
              return function () {
                var $current = 0;
                var $continue = function ($yield, $yieldin, $reject, $done) {
                  while (true) {
                    switch ($current) {
                      case 0:
                        $yield($t.fastbox(1, $g.____testlib.basictypes.Integer));
                        $current = 1;
                        return;

                      case 1:
                        $yield($t.fastbox(2, $g.____testlib.basictypes.Integer));
                        $current = 2;
                        return;

                      case 2:
                        $yield($t.fastbox(3, $g.____testlib.basictypes.Integer));
                        $current = 3;
                        return;

                      default:
                        $done();
                        return;
                    }
                  }
                };
                return $generator.new($continue);
              }().then(function ($result2) {
                return $g.children.SimpleFunction($result1, $result2).then(function ($result0) {
                  $result = $result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
