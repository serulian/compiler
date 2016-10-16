$module('classprops', function () {
  var $static = this;
  this.$class('SomeProps', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (BoolValue, StringValue) {
      var instance = new $static();
      var init = [];
      instance.BoolValue = BoolValue;
      instance.StringValue = StringValue;
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.classprops.SomeProps).$typeref()]);
    };
  });

  $static.SimpleFunction = function (props) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.String.$equals(props.StringValue, $t.box("hello world", $g.____testlib.basictypes.String)).then(function ($result2) {
              return $promise.resolve($t.unbox($result2)).then(function ($result1) {
                return $promise.resolve($result1 && $t.unbox(props.BoolValue)).then(function ($result0) {
                  $result = $t.box($result0 && !(props.OptionalValue == null), $g.____testlib.basictypes.Boolean);
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
  $static.TEST = function () {
    var $result;
    var $temp0;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.classprops.SomeProps.new(true, $t.box("hello world", $g.____testlib.basictypes.String)).then(function ($result1) {
              $temp0 = $result1;
              return $g.classprops.SimpleFunction(($temp0, $temp0.OptionalValue = $t.box(42, $g.____testlib.basictypes.Integer), $temp0)).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
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
