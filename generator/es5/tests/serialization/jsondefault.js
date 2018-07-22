$module('jsondefault', function () {
  var $static = this;
  this.$struct('d610e0b2', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
      };
      instance.$markruntimecreated();
      return $static.$initDefaults(instance, true);
    };
    $static.$initDefaults = function (instance, isRuntimeCreated) {
      var boxed = instance[BOXED_DATA_PROPERTY];
      if (isRuntimeCreated || (boxed['SomeField'] === undefined)) {
        instance.SomeField = $t.fastbox(42, $g.________testlib.basictypes.Integer);
      }
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|cf412abd<d610e0b2>": true,
        "equals|4|cf412abd<aa28dc2d>": true,
        "Stringify|2|cf412abd<cb470bcc>": true,
        "Mapping|2|cf412abd<899aec48<any>>": true,
        "Clone|2|cf412abd<d610e0b2>": true,
        "String|2|cf412abd<cb470bcc>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
    var $result;
    var jsonString;
    var parsed;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            jsonString = $t.fastbox('{}', $g.________testlib.basictypes.String);
            $promise.maybe($g.jsondefault.SomeStruct.Parse($g.________testlib.basictypes.JSON)(jsonString)).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            parsed = $result;
            $resolve($t.fastbox(parsed.SomeField.$wrapped == 42, $g.________testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
